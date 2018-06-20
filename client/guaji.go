package client

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"meatfloss/common"
	"meatfloss/gameconf"
	"meatfloss/gameuser"
	"meatfloss/message"
	"meatfloss/persistent"
	"runtime"
	"strconv"
	"time"

	"github.com/mohae/deepcopy"
)

// HandleMakingReq ...
func (c *GameClient) HandleMakingReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {
	reply := &message.MakingNotify{}
	reply.Meta.MessageType = "MakingNotify"
	reply.Meta.MessageTypeID = message.MsgMakingNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID
	req := &message.MakingReq{}
	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		return
	}

	var goodsIDsReply []string
	var goodsCountsReply []int
	var goodsNumDeltaReply []int

	//req.Data.Type==1  1添加制作   2完成制作    3解锁
	//数据表中  1制作   0空闲   3锁定
	//制作的添加
	if req.Data.Type == "1" {
		//这是个笨方法，vscode编译 if判断中的变量不能后面的代码识别
		//判断类型
		rs := []rune(req.Data.GoodsID)
		if string(rs[0:2]) == "fs" {
			//服饰
			material := gameconf.AllApparels[req.Data.GoodsID]
			num := len(material.NeedMaterial[0].List)
			for a := 0; a < num; a++ {
				//判断每件在背包中的数量
				if _, ok := c.user.Bag.Cells[gameconf.AllSuperGoods[material.NeedMaterial[0].List[a].GoodsID].UniqueID]; ok {
					//物品在背包中、
					backNum := c.user.Bag.Cells[gameconf.AllSuperGoods[material.NeedMaterial[0].List[a].GoodsID].UniqueID].Count
					_ = backNum
					if backNum < material.NeedMaterial[0].List[a].GoodsNum {
						reply.Meta.Error = true
						reply.Meta.ErrorMessage = "invalid request"
						c.SendMsg(reply)
						return
					}
				} else {
					reply.Meta.Error = true
					reply.Meta.ErrorMessage = "invalid request"
					c.SendMsg(reply)
					return
				}
			}
			//更改背包中的数量
			for a := 0; a < num; a++ {
				c.user.Bag.Cells[gameconf.AllSuperGoods[material.NeedMaterial[0].List[a].GoodsID].UniqueID].Count -= material.NeedMaterial[0].List[a].GoodsNum
				//通知背包中更改的数据
				goodsIDsReply = append(goodsIDsReply, material.NeedMaterial[0].List[a].GoodsID)
				goodsCountsReply = append(goodsCountsReply, c.user.Bag.Cells[gameconf.AllSuperGoods[material.NeedMaterial[0].List[a].GoodsID].UniqueID].Count)
				goodsNumDeltaReply = append(goodsNumDeltaReply, -material.NeedMaterial[0].List[a].GoodsNum)

			}
			c.user.MakeBox.Lattice[req.Data.Lattice-1].Required = gameconf.AllApparels[req.Data.GoodsID].Materialneed
			c.user.MakeBox.Lattice[req.Data.Lattice-1].Time = gameconf.AllApparels[req.Data.GoodsID].Maketime
			c.user.MakeBox.Lattice[req.Data.Lattice-1].End = req.Data.GoodsID
			c.user.MakeBox.Lattice[req.Data.Lattice-1].Type = 1
		} else {
			//家具
			material := gameconf.AllFurniture[req.Data.GoodsID]
			num := len(material.NeedMaterial[0].List)
			for a := 0; a < num; a++ {
				//判断每件在背包中的数量
				if _, ok := c.user.Bag.Cells[gameconf.AllSuperGoods[material.NeedMaterial[0].List[a].GoodsID].UniqueID]; ok {
					//物品在背包中、
					backNum := c.user.Bag.Cells[gameconf.AllSuperGoods[material.NeedMaterial[0].List[a].GoodsID].UniqueID].Count
					_ = backNum
					if backNum < material.NeedMaterial[0].List[a].GoodsNum {
						reply.Meta.Error = true
						reply.Meta.ErrorMessage = "invalid request"
						c.SendMsg(reply)
						return
					}
				} else {
					reply.Meta.Error = true
					reply.Meta.ErrorMessage = "invalid request"
					c.SendMsg(reply)
					return
				}
			}
			//更改背包中的数量
			for a := 0; a < num; a++ {
				c.user.Bag.Cells[gameconf.AllSuperGoods[material.NeedMaterial[0].List[a].GoodsID].UniqueID].Count -= material.NeedMaterial[0].List[a].GoodsNum
				//通知背包中更改的数据
				goodsIDsReply = append(goodsIDsReply, material.NeedMaterial[0].List[a].GoodsID)
				goodsCountsReply = append(goodsCountsReply, c.user.Bag.Cells[gameconf.AllSuperGoods[material.NeedMaterial[0].List[a].GoodsID].UniqueID].Count)
				goodsNumDeltaReply = append(goodsNumDeltaReply, -material.NeedMaterial[0].List[a].GoodsNum)
			}
			c.user.MakeBox.Lattice[req.Data.Lattice-1].Required = gameconf.AllFurniture[req.Data.GoodsID].MaterialNeed
			c.user.MakeBox.Lattice[req.Data.Lattice-1].Time = gameconf.AllFurniture[req.Data.GoodsID].MakeTime
			c.user.MakeBox.Lattice[req.Data.Lattice-1].End = req.Data.GoodsID
			c.user.MakeBox.Lattice[req.Data.Lattice-1].Type = 1

		}

		fmt.Println("我是制作")
	}
	//制作的完成
	if req.Data.Type == "2" {
		//向背包添加制作成功的物品
		if _, ok := c.user.Bag.Cells[gameconf.AllSuperGoods[req.Data.GoodsID].UniqueID]; ok {
			//物品在背包中、
			c.user.Bag.Cells[gameconf.AllSuperGoods[req.Data.GoodsID].UniqueID].Count++
			//通知背包中更改的数据
			goodsIDsReply = append(goodsIDsReply, req.Data.GoodsID)
			goodsCountsReply = append(goodsCountsReply, c.user.Bag.Cells[gameconf.AllSuperGoods[req.Data.GoodsID].UniqueID].Count)
			goodsNumDeltaReply = append(goodsNumDeltaReply, 1)
		} else {
			//物品没有在背包中
			var goodsIDs []string
			var goodsCounts []int
			goodsIDs = append(goodsIDs, req.Data.GoodsID)
			goodsCounts = append(goodsCounts, 1)
			c.PutToBagBatch(goodsIDs, goodsCounts)
			//通知背包中更改的数据
			goodsIDsReply = append(goodsIDsReply, req.Data.GoodsID)
			goodsCountsReply = append(goodsCountsReply, 1)
			goodsNumDeltaReply = append(goodsNumDeltaReply, 1)
		}
		//更改格子的状态
		c.user.MakeBox.Lattice[req.Data.Lattice-1].Required = "0"
		c.user.MakeBox.Lattice[req.Data.Lattice-1].Time = 0
		c.user.MakeBox.Lattice[req.Data.Lattice-1].End = "0"
		c.user.MakeBox.Lattice[req.Data.Lattice-1].Type = 0
	}
	//解锁
	if req.Data.Type == "3" {
		//只能是被锁定的情况下才能被解锁
		if c.user.MakeBox.Lattice[req.Data.Lattice-1].Type == 3 {
			//判断钻石数
			if _, ok := c.user.Bag.Cells[gameconf.AllSuperGoods["wp0001"].UniqueID]; ok {
				//物品在背包中、
				backNum := c.user.Bag.Cells[gameconf.AllSuperGoods["wp0001"].UniqueID].Count
				_ = backNum
				if backNum < gameconf.AllLattice[req.Data.Lattice].UnlockPrice {
					reply.Meta.Error = true
					reply.Meta.ErrorMessage = "invalid request"
					c.SendMsg(reply)
					return
				}
			} else {
				//钻石在背包中也没有
				reply.Meta.Error = true
				reply.Meta.ErrorMessage = "invalid request"
				c.SendMsg(reply)
				return
			}
			//更新背包中钻石的数量
			c.user.Bag.Cells[gameconf.AllSuperGoods["wp0001"].UniqueID].Count -= gameconf.AllLattice[req.Data.Lattice].UnlockPrice
			c.user.MakeBox.Lattice[req.Data.Lattice-1].Type = 0
		}
	}
	//加速
	if req.Data.Type == "4" {
		//先判断背包中的数量是否满足加速的需求
		if _, ok := c.user.Bag.Cells[gameconf.AllSuperGoods["wp0001"].UniqueID]; ok {
			//物品在背包中、
			backNum := c.user.Bag.Cells[gameconf.AllSuperGoods["wp0001"].UniqueID].Count
			_ = backNum
			if backNum < gameconf.AllLattice[req.Data.Lattice].UnlockPrice {
				reply.Meta.Error = true
				reply.Meta.ErrorMessage = "invalid request"
				c.SendMsg(reply)
				return
			}
		} else {
			//钻石在背包中也没有
			reply.Meta.Error = true
			reply.Meta.ErrorMessage = "invalid request"
			c.SendMsg(reply)
			return
		}
		//更新背包的数量
		c.user.Bag.Cells[gameconf.AllSuperGoods["wp0001"].UniqueID].Count -= req.Data.CoinNum
		//更新格子的时间为0
		c.user.MakeBox.Lattice[req.Data.Lattice-1].Time = 0
	}

	//拿到服饰
	temporaryApparel := make(map[string]string)
	temporary := &message.MakeLatticeBack{}
	_ = temporary
	apparel := gameconf.AllApparels
	num := len(apparel)
	for j := 1; j <= num; j++ {
		key := ""
		if j <= 9 {
			key = "fs000"
		} else {
			key = "fs00"
		}
		string := strconv.Itoa(j)
		key += string

		if _, ok := c.user.Bag.Cells[gameconf.AllSuperGoods[apparel[key].ID].UniqueID]; ok {
			//此物品在背包中存在
			backNum := c.user.Bag.Cells[gameconf.AllSuperGoods[apparel[key].ID].UniqueID].Count
			_ = backNum
			string := strconv.Itoa(backNum)
			numa := len(c.user.CollectionBox.Collections)
			base := 0
			if numa > 0 {
				for a := 0; a < numa; a++ {
					if c.user.CollectionBox.Collections[a].GoodID == apparel[key].ID {
						base = 1
					}
				}
			}
			if base == 1 {
				temporaryApparel[apparel[key].ImageName] += apparel[key].ID + "|" + apparel[key].ImageName + "|" + string + "|1;"
			} else {
				temporaryApparel[apparel[key].ImageName] += apparel[key].ID + "|" + apparel[key].ImageName + "|" + string + "|0;"
			}
		} else {
			//背包中没有

			num := len(c.user.CollectionBox.Collections)
			base := 0
			if num > 0 {
				for a := 0; a < num; a++ {
					if c.user.CollectionBox.Collections[a].GoodID == apparel[key].ID {
						base = 1
					}
				}
			}
			if base == 1 {
				temporaryApparel[apparel[key].ImageName] += apparel[key].ID + "|" + apparel[key].ImageName + "|0|1;"
			} else {
				temporaryApparel[apparel[key].ImageName] += apparel[key].ID + "|" + apparel[key].ImageName + "|0|0;"
			}

		}
	}
	//拿到家具
	temporaryFuniture := make(map[string]string)

	funiture := gameconf.AllFurniture
	numFun := len(funiture)
	for j := 1; j <= numFun; j++ {
		key := ""
		if j <= 9 {
			key = "jj000"
		} else {
			key = "jj00"
		}
		string := strconv.Itoa(j)
		key += string
		fmt.Println(funiture[key].UniqueID)
		if _, ok := c.user.Bag.Cells[funiture[key].UniqueID]; ok {

			backNum := c.user.Bag.Cells[funiture[key].UniqueID].Count
			_ = backNum
			string := strconv.Itoa(backNum)

			numa := len(c.user.CollectionBox.Collections)
			base := 0
			if numa > 0 {
				for a := 0; a < numa; a++ {
					if c.user.CollectionBox.Collections[a].GoodID == funiture[key].ID {
						base = 1
					}
				}
			}
			if base == 1 {
				temporaryFuniture[funiture[key].ImageName] += funiture[key].ID + "|" + funiture[key].ImageName + "|" + string + "|1;"
			} else {
				temporaryFuniture[funiture[key].ImageName] += funiture[key].ID + "|" + funiture[key].ImageName + "|" + string + "|0;"
			}
		} else {
			numa := len(c.user.CollectionBox.Collections)
			base := 0
			if numa > 0 {
				for a := 0; a < numa; a++ {
					if c.user.CollectionBox.Collections[a].GoodID == funiture[key].ID {
						base = 1
					}
				}
			}
			if base == 1 {
				temporaryFuniture[funiture[key].ImageName] += funiture[key].ID + "|" + funiture[key].ImageName + "|0|1;"
			} else {
				temporaryFuniture[funiture[key].ImageName] += funiture[key].ID + "|" + funiture[key].ImageName + "|0|0;"
			}
		}
	}
	temporary.Furniture = temporaryFuniture
	temporary.Clothes = temporaryApparel
	{
		cpy := deepcopy.Copy(temporary)
		layout, _ := cpy.(*message.MakeLatticeBack)
		reply.Data.MakeLatticeBack = layout
	}

	//信息的发送
	{
		cpy := deepcopy.Copy(c.user.MakeBox)
		mak, _ := cpy.(*gameuser.MakeBox)
		reply.Data.Lattice = mak.Lattice
	}
	//收藏的推送
	{
		cpy := deepcopy.Copy(c.user.CollectionBox)
		collection, _ := cpy.(*gameuser.CollectionBox)
		num := len(collection.Collections)

		if num > 0 {
			for a := 0; a < num; a++ {
				collect := &common.Collections{}
				if _, ok := c.user.Bag.Cells[gameconf.AllSuperGoods[collection.Collections[a].GoodID].UniqueID]; ok {
					collect.GoodsNum = c.user.Bag.Cells[gameconf.AllSuperGoods[collection.Collections[a].GoodID].UniqueID].Count

				} else {
					collect.GoodsNum = 0

				}
				collect.GoodID = collection.Collections[a].GoodID
				reply.Data.Collection = append(reply.Data.Collection, collect)
			}
		} else {
			reply.Data.Collection = collection.Collections
		}

	}

	notify := message.UpdateGoodsNotify{}
	notify.Meta.MessageType = "UpdateGoodsNotify"
	notify.Meta.MessageTypeID = message.MsgTypeUpdateGoodsNotify

	//返回背包中更改的数据
	updateNum := len(goodsCountsReply)
	_ = updateNum
	if updateNum > 0 {
		for a := 0; a < updateNum; a++ {
			//查询当前的物品是否可堆叠
			if c.GetType(goodsIDsReply[a]) == 1 { //可堆叠
				update := &message.GoodsUpdateInfo{}
				update.GoodsID = goodsIDsReply[a]
				update.GoodsNum = goodsCountsReply[a]
				update.GoodsNumDelta = goodsNumDeltaReply[a]
				update.UniqueID = gameconf.AllSuperGoods[goodsIDsReply[a]].UniqueID
				notify.Data.List = append(notify.Data.List, *update)
			} else { //不可堆叠
				update := &message.GoodsUpdateInfo{}
				update.GoodsID = goodsIDsReply[a]
				update.GoodsNum = goodsCountsReply[a]
				update.GoodsNumDelta = goodsNumDeltaReply[a]
				update.UniqueID = gameconf.AllSuperGoods[goodsIDsReply[a]].UniqueID*10 + int64(goodsCountsReply[a])
				notify.Data.List = append(notify.Data.List, *update)
			}
		}
		notify.Data.Type = "2"
		c.SendMsg(notify)
	}
	c.PushUserNotify()
	c.persistMaking()
	c.persistBagBox()
	c.SendMsg(reply)
	return
}
func (c *GameClient) persistMaking() {

	making := &gameuser.User{}

	cpy := deepcopy.Copy(c.user.MakeBox)
	mak, _ := cpy.(*gameuser.MakeBox)
	making.MakeBox = mak
	persistent.AddUser(c.UserID, making)

}

// HandleCollectionReq ...
func (c *GameClient) HandleCollectionReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	reply := &message.CollectionNotify{}
	reply.Meta.MessageType = "CollectionNotify"
	reply.Meta.MessageTypeID = message.MsgCollectionNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID
	req := &message.CollectionReq{}
	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		return
	}
	//req.Data.Type   type 1添加收藏
	Collection := c.user.CollectionBox
	//先看是否已经书藏过了
	num := len(Collection.Collections)
	_ = num
	base := 1
	del := 0
	_ = base
	_ = del
	if num > 0 {
		for a := 0; a < num; a++ {
			if Collection.Collections[a].GoodID == req.Data.GoodsID {
				base = 2
				del = a
			}
		}
	}
	//添加收藏收藏
	if req.Data.Type == "1" {
		//如果base为1就添加收藏
		if base == 1 {
			Collection := &common.Collections{}
			Collection.GoodID = req.Data.GoodsID
			c.user.CollectionBox.Collections = append(c.user.CollectionBox.Collections, Collection)
		}
	} else {
		//取消收藏
		if base == 2 {
			if num > 0 {
				for a := 0; a < num; a++ {
					if del == a {
						Collection.Collections = append(Collection.Collections[:a], Collection.Collections[a+1:]...)
					}
				}
			}
		}
	}
	//收藏的通知
	{
		cpy := deepcopy.Copy(c.user.CollectionBox)
		collection, _ := cpy.(*gameuser.CollectionBox)
		num := len(collection.Collections)
		if num > 0 {
			for a := 0; a < num; a++ {
				collect := &common.Collections{}
				//判断背包中是否有这个收藏的物品
				if _, ok := c.user.Bag.Cells[gameconf.AllSuperGoods[collection.Collections[a].GoodID].UniqueID]; ok {
					collect.GoodsNum = c.user.Bag.Cells[gameconf.AllSuperGoods[collection.Collections[a].GoodID].UniqueID].Count
				} else {
					collect.GoodsNum = 0
				}
				collect.GoodID = collection.Collections[a].GoodID
				reply.Data.Collection = append(reply.Data.Collection, collect)
			}
		} else {
			reply.Data.Collection = collection.Collections
		}

	}

	c.SendMsg(reply)
	c.persistCollection()
	return
}

func (c *GameClient) persistCollection() {

	Collection := &gameuser.User{}

	cpy := deepcopy.Copy(c.user.CollectionBox)
	collection, _ := cpy.(*gameuser.CollectionBox)
	Collection.CollectionBox = collection
	persistent.AddUser(c.UserID, Collection)

}

// HandleMakeLatticeReq ...
func (c *GameClient) HandleMakeLatticeReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	reply := &message.MakeLatticNotify{}
	reply.Meta.MessageType = "MakeLatticNotify"
	reply.Meta.MessageTypeID = message.MsgMakeLatticeNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID
	req := &message.MakeLatticNotify{}
	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		return
	}

	//拿到服饰
	temporaryApparel := make(map[string]string)
	temporary := &message.MakeLatticeBack{}
	_ = temporary
	apparel := gameconf.AllApparels
	num := len(apparel)
	for j := 1; j <= num; j++ {
		key := ""
		if j <= 9 {
			key = "fs000"
		} else {
			key = "fs00"
		}
		string := strconv.Itoa(j)
		key += string

		if _, ok := c.user.Bag.Cells[gameconf.AllSuperGoods[apparel[key].ID].UniqueID]; ok {
			//此物品在背包中存在
			backNum := c.user.Bag.Cells[gameconf.AllSuperGoods[apparel[key].ID].UniqueID].Count
			_ = backNum
			string := strconv.Itoa(backNum)
			numa := len(c.user.CollectionBox.Collections)
			base := 0
			if numa > 0 {
				for a := 0; a < numa; a++ {
					if c.user.CollectionBox.Collections[a].GoodID == apparel[key].ID {
						base = 1
					}
				}
			}
			if base == 1 {
				temporaryApparel[apparel[key].ImageName] += apparel[key].ID + "|" + apparel[key].ImageName + "|" + string + "|1;"
			} else {
				temporaryApparel[apparel[key].ImageName] += apparel[key].ID + "|" + apparel[key].ImageName + "|" + string + "|0;"
			}
		} else {
			//背包中没有

			num := len(c.user.CollectionBox.Collections)
			base := 0
			if num > 0 {
				for a := 0; a < num; a++ {
					if c.user.CollectionBox.Collections[a].GoodID == apparel[key].ID {
						base = 1
					}
				}
			}
			if base == 1 {
				temporaryApparel[apparel[key].ImageName] += apparel[key].ID + "|" + apparel[key].ImageName + "|0|1;"
			} else {
				temporaryApparel[apparel[key].ImageName] += apparel[key].ID + "|" + apparel[key].ImageName + "|0|0;"
			}

		}
	}
	//拿到家具
	temporaryFuniture := make(map[string]string)

	funiture := gameconf.AllFurniture
	numFun := len(funiture)
	for j := 1; j <= numFun; j++ {
		key := ""
		if j <= 9 {
			key = "jj000"
		} else {
			key = "jj00"
		}
		string := strconv.Itoa(j)
		key += string
		fmt.Println(funiture[key].UniqueID)
		if _, ok := c.user.Bag.Cells[funiture[key].UniqueID]; ok {

			backNum := c.user.Bag.Cells[funiture[key].UniqueID].Count
			_ = backNum
			string := strconv.Itoa(backNum)

			numa := len(c.user.CollectionBox.Collections)
			base := 0
			if numa > 0 {
				for a := 0; a < numa; a++ {
					if c.user.CollectionBox.Collections[a].GoodID == funiture[key].ID {
						base = 1
					}
				}
			}
			if base == 1 {
				temporaryFuniture[funiture[key].ImageName] += funiture[key].ID + "|" + funiture[key].ImageName + "|" + string + "|1;"
			} else {
				temporaryFuniture[funiture[key].ImageName] += funiture[key].ID + "|" + funiture[key].ImageName + "|" + string + "|0;"
			}
		} else {
			numa := len(c.user.CollectionBox.Collections)
			base := 0
			if numa > 0 {
				for a := 0; a < numa; a++ {
					if c.user.CollectionBox.Collections[a].GoodID == funiture[key].ID {
						base = 1
					}
				}
			}
			if base == 1 {
				temporaryFuniture[funiture[key].ImageName] += funiture[key].ID + "|" + funiture[key].ImageName + "|0|1;"
			} else {
				temporaryFuniture[funiture[key].ImageName] += funiture[key].ID + "|" + funiture[key].ImageName + "|0|0;"
			}
		}
	}
	temporary.Furniture = temporaryFuniture
	temporary.Clothes = temporaryApparel
	{
		cpy := deepcopy.Copy(temporary)
		layout, _ := cpy.(*message.MakeLatticeBack)
		reply.Data.MakeLatticeBack = layout
	}
	//收藏的通知
	{
		cpy := deepcopy.Copy(c.user.CollectionBox)
		collection, _ := cpy.(*gameuser.CollectionBox)
		num := len(collection.Collections)

		if num > 0 {
			for a := 0; a < num; a++ {
				collect := &common.Collections{}
				//判断背包中是否有这个收藏的物品
				if _, ok := c.user.Bag.Cells[gameconf.AllSuperGoods[collection.Collections[a].GoodID].UniqueID]; ok {
					collect.GoodsNum = c.user.Bag.Cells[gameconf.AllSuperGoods[collection.Collections[a].GoodID].UniqueID].Count

				} else {
					collect.GoodsNum = 0
				}
				collect.GoodID = collection.Collections[a].GoodID
				reply.Data.Collection = append(reply.Data.Collection, collect)
			}
		} else {
			reply.Data.Collection = collection.Collections
		}
	}

	//制作的通知
	{
		cpy := deepcopy.Copy(c.user.MakeBox)
		mak, _ := cpy.(*gameuser.MakeBox)
		reply.Data.Lattice = mak.Lattice
	}
	c.SendMsg(reply)
	c.PushUserNotify()
	return
}

// HandleWgReq ...
func (c *GameClient) HandleWgReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {
	reply := &message.WgNotify{}
	reply.Meta.MessageType = "WgNotify"
	reply.Meta.MessageTypeID = message.MsgWgtNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID
	req := &message.ShowWgReq{}
	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		return
	}
	inta, err := strconv.Atoi(req.Data.Num)
	userID, err := strconv.Atoi(req.Data.UserID)
	_ = userID
	//goodsID, err := strconv.Atoi(req.Data.GoodsID)

	var goodsIDs []string
	var goodsCounts []int
	goodsIDs = append(goodsIDs, req.Data.GoodsID)
	goodsCounts = append(goodsCounts, inta)
	c.PutToBagBatch(goodsIDs, goodsCounts)

	newUser := &gameuser.User{}
	newUser.UserID = c.UserID

	cpy := deepcopy.Copy(c.user.Bag)
	bag, _ := cpy.(*common.Bag)
	newUser.Bag = bag
	persistent.AddUser(c.UserID, newUser)
	reply.Data.Res = "success"
	c.SendMsg(reply)
	c.PushUserNotify()
	//persistent.AddUser(userID, newGuajiProfile)
	return
}

// HandleMachineUpgradeReq ...
func (c *GameClient) HandleMachineUpgradeReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	reply := &message.MachineUpgradeNotify{}
	reply.Meta.MessageType = "MachineUpgradeNotify"
	reply.Meta.MessageTypeID = message.MsgMyUpgradeNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID
	req := &message.ShowMachineUpgradeReq{}

	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		return
	}
	// req.Data.MachineUpgradeApply   1就是升级的请求
	// req.Data.MachineUpgradeConfirm 1就是确认升级的请求
	//处理请求升级的请求
	if req.Data.MachineUpgradeApply == "1" && req.Data.MachineUpgradeConfirm == "2" {
		//返还当当前的等级的数据还有升级之后的等级的数据
		myEmployee := &message.RoleGuajiSettlement{}
		myEmployee.Luck = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Luck
		myEmployee.Quality = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Quality
		myEmployee.Speed = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Speed
		myEmployee.NumEmployees = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].NumEmployees
		myEmployee.MachineLevel = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].MachineLevel
		myEmployee.MachineImage = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].MachineImage
		myEmployee.CD = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].CD
		myEmployee.CDPerDegree = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].CDPerDegree
		myEmployee.Uptime = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Uptime
		myEmployee.MaxTemperature = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].MaxTemperature
		myEmployee.Upmaterial = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Upmaterial
		// machineNeed := gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Guajis[0].List
		// //判断升级需要的材料
		// var material string
		// if len(machineNeed) >= 0 {
		// 	//循环升级的材料，查出他的信息
		// 	for goodsId, goods := range machineNeed {
		// 		_ = goodsId
		// 		GoodsNum := strconv.Itoa(goods.GoodsNum)
		// 		material += gameconf.AllGoods[goods.GoodsID].Name + ":" + GoodsNum + ";"
		// 	}
		// }
		// myEmployee.Upmaterial = material
		reply.Data.MachineNow = append(reply.Data.MachineNow, myEmployee)
		//返还升级后的
		myEmployeeUp := &message.RoleGuajiSettlement{}
		myEmployeeUp.Luck = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].Luck
		myEmployeeUp.Quality = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].Quality
		myEmployeeUp.Speed = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].Speed
		myEmployeeUp.NumEmployees = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].NumEmployees
		myEmployeeUp.MachineLevel = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].MachineLevel
		myEmployeeUp.MachineImage = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].MachineImage
		myEmployeeUp.CD = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].CD
		myEmployeeUp.CDPerDegree = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].CDPerDegree
		myEmployeeUp.Uptime = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].Uptime
		myEmployeeUp.MaxTemperature = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].MaxTemperature
		myEmployeeUp.Upmaterial = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].Upmaterial
		// machineNeedUp := gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel+1].Guajis[0].List
		// //判断升级需要的材料
		// var materialUp string
		// if len(machineNeedUp) >= 0 {
		// 	//循环升级的材料，查出他的信息
		// 	for goodsId, goods := range machineNeedUp {
		// 		_ = goodsId
		// 		GoodsNum := strconv.Itoa(goods.GoodsNum)
		// 		materialUp += gameconf.AllGoods[goods.GoodsID].Name + ":" + GoodsNum + ";"
		// 	}
		// }
		// myEmployeeUp.Upmaterial = materialUp
		reply.Data.MachineUpgrade = append(reply.Data.MachineUpgrade, myEmployeeUp)
		c.SendMsg(reply)
		return
	}
	//确认升级的
	if req.Data.MachineUpgradeApply == "2" && req.Data.MachineUpgradeConfirm == "1" {
		//判断是否达到升级的要求
		machineNeedConfirm := gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Guajis[0].List
		//判断升级需要的材料
		var Whether int32 = 0
		_ = Whether
		var materialNeed string
		if len(machineNeedConfirm) >= 0 {
			//循环在背包比较需要的材料，背包中的数量是否满足
			for goodsId, goods := range machineNeedConfirm {
				_ = goodsId
				if _, ok := c.user.Bag.Cells[gameconf.AllSuperGoods[goods.GoodsID].UniqueID]; ok {
					//此物品在背包中存在
					GoodsNum := strconv.Itoa(goods.GoodsNum)
					if c.user.Bag.Cells[gameconf.AllSuperGoods[goods.GoodsID].UniqueID].Count < goods.GoodsNum {
						materialNeed += gameconf.AllGoods[goods.GoodsID].Name + "数量不足：需要" + GoodsNum + "个！"
						Whether = 1
					}
				} else {
					//物品在背包中不存在
					GoodsNum := strconv.Itoa(goods.GoodsNum)
					//if c.user.Bag.Cells[gameconf.AllSuperGoods[goods.GoodsID].UniqueID].Count < goods.GoodsNum {
					materialNeed += gameconf.AllGoods[goods.GoodsID].Name + "数量不足：需要" + GoodsNum + "个！"
					Whether = 1
				}

			}
		}
		if Whether == 1 {

			reply.Data.MachineUpgradeType = "no"
			reply.Data.MachineUpgradeMessage = materialNeed
			c.SendMsg(reply)
			return
		} else {
			if c.user.GuajiProfile.Upgrade != 1 {
				c.user.GuajiProfile.Upgrade = 2
			}
			reply.Data.MachineUpgradeType = "yes"
			reply.Data.MachineUpgradeMessage = "正在升级"
			c.SendMsg(reply)
			c.persistGuajiProfile()
			return
		}
	}

	reply.Meta.Error = true
	reply.Meta.ErrorMessage = "invalid request"
	c.SendMsg(reply)
	return

}

// HandleMyEmployeeReq ...
func (c *GameClient) HandleMyEmployeeReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	reply := &message.MyEmployeeNotify{}
	reply.Meta.MessageType = "MyEmployeeNotify"
	reply.Meta.MessageTypeID = message.MsgMyEmployeeNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID
	var num int = len(c.user.GuajiProfile.EmployeeBox.EmployeesInfo)
	var numB int = len(c.user.Bag.BagEmployee)
	if num == 0 && numB == 0 {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "invalid request"
		c.SendMsg(reply)
		return
	}
	if num != 0 {
		//工作的
		for a := 0; a < num; a++ {
			//	go func(who int) {
			myEmployee := &message.Employeeinfo{}
			numid := c.user.GuajiProfile.EmployeeBox.EmployeesInfo[a].EmployeesID
			myEmployee.Speed = gameconf.AllEmployees[numid].Speed
			myEmployee.Quality = gameconf.AllEmployees[numid].Quality
			myEmployee.Number = gameconf.AllEmployees[numid].Number
			myEmployee.Luck = gameconf.AllEmployees[numid].Luck
			myEmployee.Introdution = gameconf.AllEmployees[numid].Introdution
			myEmployee.EmployeeName = gameconf.AllEmployees[numid].EmployeeName
			myEmployee.AvatarImage = gameconf.AllEmployees[numid].AvatarImage
			reply.Data.EmployeeWork = append(reply.Data.EmployeeWork, myEmployee)
			//}(a)
		}
	}
	if numB != 0 {
		//背包
		for a := 0; a < numB; a++ {
			//	go func(who int) {
			myEmployee := &message.Employeeinfo{}
			numid := c.user.Bag.BagEmployee[a].EmployeesID
			myEmployee.Speed = gameconf.AllEmployees[numid].Speed
			myEmployee.Quality = gameconf.AllEmployees[numid].Quality
			myEmployee.Number = gameconf.AllEmployees[numid].Number
			myEmployee.Luck = gameconf.AllEmployees[numid].Luck
			myEmployee.Introdution = gameconf.AllEmployees[numid].Introdution
			myEmployee.EmployeeName = gameconf.AllEmployees[numid].EmployeeName
			myEmployee.AvatarImage = gameconf.AllEmployees[numid].AvatarImage
			reply.Data.EmployeeBack = append(reply.Data.EmployeeBack, myEmployee)
			//}(a)
		}
	}
	//读取当前机器的属性值
	myEmployee := &message.RoleGuajiSettlement{}
	myEmployee.Luck = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Luck
	myEmployee.Quality = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Quality
	myEmployee.Speed = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Speed
	myEmployee.NumEmployees = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].NumEmployees
	myEmployee.MachineLevel = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].MachineLevel
	myEmployee.MachineImage = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].MachineImage
	reply.Data.Machine = append(reply.Data.Machine, myEmployee)
	c.SendMsg(reply)
	return
}

// HandleEmployeeAdjustReq ...
func (c *GameClient) HandleEmployeeAdjustReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	reply := &message.EmployeeAdjustNotify{}
	reply.Meta.MessageType = "EmployeeAdjustNotify"
	reply.Meta.MessageTypeID = message.MsgEmployeeAdjustNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID

	req := &message.SaveEmployeeAdjustReq{}

	fmt.Println(req.Data.EmployeeAdjust)

	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "invalid request"
		c.SendMsg(reply)
		return
	}
	//读取当前机器的属性值
	myEmployee := &message.RoleGuajiSettlement{}
	myEmployee.Luck = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Luck
	myEmployee.Quality = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Quality
	myEmployee.Speed = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].Speed
	myEmployee.NumEmployees = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].NumEmployees
	myEmployee.MinLevel = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].MinLevel
	myEmployee.MachineImage = gameconf.AllGuajis[c.user.GuajiProfile.MachineLevel].MachineImage
	reply.Data.Machine = append(reply.Data.Machine, myEmployee)

	//产生随机的标示值
	r := rand.New(rand.NewSource(time.Now().Unix()))
	fmt.Println(r.Intn(10000)) // [0,100)的随机值，返回值为int
	string := strconv.Itoa(r.Intn(10000))
	str := "token" + string

	cpy := deepcopy.Copy(req.Data.EmployeeAdjust)
	layout, _ := cpy.(*message.EmployeeAdjust)

	//加入工作中
	if len(layout.Employee) > 0 {
		c.user.GuajiProfile.EmployeeBox.EmployeesInfo = make([]*common.EmployeesInfo, 0)
		c.user.GuajiProfile.EmployeeBox.EmployeesToken = str
		for a := len(layout.Employee); a >= 1; a-- {
			go func(who int) {
				onet := &common.EmployeesInfo{}
				onet.EmployeesID = layout.Employee[who]
				c.user.GuajiProfile.EmployeeBox.EmployeesInfo = append(c.user.GuajiProfile.EmployeeBox.EmployeesInfo, onet)

				time.Sleep(10 * time.Nanosecond)
			}(a)
		}
		runtime.Gosched()
	} else {
		c.user.GuajiProfile.EmployeeBox.EmployeesInfo = make([]*common.EmployeesInfo, 0)
	}

	//加入背包
	if len(layout.Back) > 0 {
		c.user.Bag.BagEmployee = make([]*common.EmployeesInfo, 0)

		for a := len(layout.Back); a >= 1; a-- {
			go func(who int) {
				onet := &common.EmployeesInfo{}
				onet.EmployeesID = layout.Back[who]
				c.user.Bag.BagEmployee = append(c.user.Bag.BagEmployee, onet)

				time.Sleep(10 * time.Nanosecond)
			}(a)
		}

		runtime.Gosched()
	} else {
		c.user.Bag.BagEmployee = make([]*common.EmployeesInfo, 0)
	}
	fmt.Println(len(layout.Employee))

	var num int = len(c.user.GuajiProfile.EmployeeBox.EmployeesInfo)
	var numB int = len(c.user.Bag.BagEmployee)

	if num == 0 && numB == 0 {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "invalid request"
		c.SendMsg(reply)
		return
	}
	if num != 0 {
		//工作的
		for a := 0; a < num; a++ {
			//	go func(who int) {
			myEmployee := &message.Employeeinfo{}
			numid := c.user.GuajiProfile.EmployeeBox.EmployeesInfo[a].EmployeesID
			myEmployee.Speed = gameconf.AllEmployees[numid].Speed
			myEmployee.Quality = gameconf.AllEmployees[numid].Quality
			myEmployee.Number = gameconf.AllEmployees[numid].Number
			myEmployee.Luck = gameconf.AllEmployees[numid].Luck
			myEmployee.Introdution = gameconf.AllEmployees[numid].Introdution
			myEmployee.EmployeeName = gameconf.AllEmployees[numid].EmployeeName
			myEmployee.AvatarImage = gameconf.AllEmployees[numid].AvatarImage
			reply.Data.EmployeeWork = append(reply.Data.EmployeeWork, myEmployee)
			//}(a)
		}
	}
	if numB != 0 {
		//背包
		for a := 0; a < numB; a++ {
			//	go func(who int) {
			myEmployee := &message.Employeeinfo{}
			numid := c.user.Bag.BagEmployee[a].EmployeesID
			myEmployee.Speed = gameconf.AllEmployees[numid].Speed
			myEmployee.Quality = gameconf.AllEmployees[numid].Quality
			myEmployee.Number = gameconf.AllEmployees[numid].Number
			myEmployee.Luck = gameconf.AllEmployees[numid].Luck
			myEmployee.Introdution = gameconf.AllEmployees[numid].Introdution
			myEmployee.EmployeeName = gameconf.AllEmployees[numid].EmployeeName
			myEmployee.AvatarImage = gameconf.AllEmployees[numid].AvatarImage
			reply.Data.EmployeeBack = append(reply.Data.EmployeeBack, myEmployee)
			//}(a)
		}
	}

	c.SendMsg(reply)
	c.persistEmployee()
	c.persistBagBox()
	return
}

func (c *GameClient) persistEmployee() {

	Employee := &gameuser.User{}

	cpy := deepcopy.Copy(c.user.GuajiProfile)
	adjust, _ := cpy.(*gameuser.GuajiProfile)
	Employee.GuajiProfile = adjust
	persistent.AddUser(c.UserID, Employee)

}

// HandleEmployeeListReq ...
func (c *GameClient) HandleEmployeeListReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	reply := &message.EmployeeListNotify{}
	reply.Meta.MessageType = "EmployeeListNotify"
	reply.Meta.MessageTypeID = message.MsgEmployeeListNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID

	reply.Data.Employee = make([]*message.Employeeinfo, 0)
	fmt.Println(len(c.user.GuajiOutputBox.GuajiOutputs))

	for a := 1; a <= 10; a++ {
		//	go func(who int) {
		myEmployee := &message.Employeeinfo{}
		var str = "gy00"
		if a <= 9 {
			str = "gy00"
		} else {
			str = "gy0"
		}
		d := strconv.Itoa(a)
		str += d
		myEmployee.Speed = gameconf.AllEmployees[str].Speed
		myEmployee.Quality = gameconf.AllEmployees[str].Quality
		myEmployee.Number = gameconf.AllEmployees[str].Number
		myEmployee.Luck = gameconf.AllEmployees[str].Luck
		myEmployee.Introdution = gameconf.AllEmployees[str].Introdution
		myEmployee.EmployeeName = gameconf.AllEmployees[str].EmployeeName
		myEmployee.AvatarImage = gameconf.AllEmployees[str].AvatarImage
		reply.Data.Employee = append(reply.Data.Employee, myEmployee)
		//}(a)
	}
	c.SendMsg(reply)
	return
}

//回复前端的信息
//HandleOutputReq  ...
func (c *GameClient) HandleOutputReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	reply := &message.OutputNotify{}
	reply.Meta.MessageType = "OutputNotify"
	reply.Meta.MessageTypeID = message.MsgTypeOutputNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID

	if len(c.user.GuajiOutputBox.GuajiOutputs) == 0 {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "GuajiOutputs don't exits"
		c.SendMsg(reply)
		return
	}
	//消息的推送
	fmt.Println(c.user.GuajiOutputBox.GuajiOutputs)
	for i := len(c.user.GuajiOutputBox.GuajiOutputs) - 1; i >= 0; i-- {
		guajiOutput := &common.GuajiOutputInfo{}
		guajiOutput.UserID = c.user.GuajiOutputBox.GuajiOutputs[i].UserID
		guajiOutput.Type = c.user.GuajiOutputBox.GuajiOutputs[i].Type
		guajiOutput.Name = c.user.GuajiOutputBox.GuajiOutputs[i].Name
		guajiOutput.Time = string([]byte(c.user.GuajiOutputBox.GuajiOutputs[i].Time)[:16])
		guajiOutput.Items = c.user.GuajiOutputBox.GuajiOutputs[i].Items
		reply.Data.GuajiOutputs = append(reply.Data.GuajiOutputs, guajiOutput)
	}
	c.SendMsg(reply)
	return
}

func (c *GameClient) persistOutput() {
	newOutput := &gameuser.User{}
	newOutput.UserID = c.UserID
	cpy := deepcopy.Copy(c.user.GuajiOutputBox)
	output, _ := cpy.(*gameuser.GuajiOutputBox)
	newOutput.GuajiOutputBox = output
	persistent.AddUser(c.UserID, newOutput)
}
