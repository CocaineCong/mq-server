package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mq-server/cache"
	"mq-server/dao"
	"mq-server/model"
	"strconv"
)

func MQ2MySQL() {
	r := cache.RedisClient      // redis
	ch, _ := model.MQ.Channel() //打开Channel
	q, _ := ch.QueueDeclare("skill_goods", true, false, false, false, nil)
	_ = ch.Qos(1, 0, false)
	msgs, _ := ch.Consume(q.Name, "", false, false, false, false, nil)
	orderDao := dao.NewOrderDao(context.Background())
	go func() {
		for d := range msgs { // 开始消费
			var p model.SkillGood2MQ
			_ = json.Unmarshal(d.Body, &p)
			// 创建订单
			order := &model.Order{
				UserID:    p.UserId,
				ProductID: p.ProductId,
				BossID:    p.BossId,
				AddressID: p.AddressId,
				Num:       1,
				OrderNum:  1,
				Type:      0,
				Money:     p.Money,
			}
			errOrder := orderDao.CreateOrder(order)
			if errOrder != nil {
				return
			}
			// 支付
			paymentObj := &OrderPay{
				OrderId:   order.ID,
				Money:     order.Money * float64(order.Num),
				ProductID: order.ProductID,
				BossID:    order.BossID,
				Num:       order.Num,
				Key:       p.Key,
			}
			fmt.Println("payment", *paymentObj)
			errPay := paymentObj.PayDown(context.Background(), p.UserId)
			if errPay != nil {
				return
			}
			// redis扣除
			r.HIncrBy("SK"+strconv.Itoa(int(p.SkillGoodId)), "num", -1) // 数量 -1

			// 存入数据库
			log.Printf("Done")
			_ = d.Ack(false) // 确认消息,必须为false
		}
	}()
}

func UpdateMySQLInfo(pId uint) {
	// 更新商品数量,直接读取redis中的结果进行更新
	n := cache.RedisClient.HGet(strconv.Itoa(int(pId)), "num").String()
	num, _ := strconv.Atoi(n)
	err := dao.NewSkillGoodsDao(context.Background()).UpdateNumByRedis(pId, num)
	if err != nil {
		return
	}
}
