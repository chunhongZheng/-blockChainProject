package main

import (
	"crypto/sha256"
)

//默克尔树节点
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

//默克尔树根结点
type MerkleTree struct {
	RootNode *MerkleNode
}

// 生成默克尔树中的节点，如果是叶子节点，则Left ,Right为nil, 如果为非叶子节点，根据Left ,Right生成当前节点的hash
func NewMerkleNode(Left, Right *MerkleNode, data []byte) *MerkleNode {
	mnode := MerkleNode{}
	if Left == nil && Right == nil {
		//当前节点为叶子节点
		mnode.Data = data
	} else {
		prehashes := append(Left.Data, Right.Data...)
		firsthash := sha256.Sum256(prehashes)
		hash := sha256.Sum256(firsthash[:])
		mnode.Data = hash[:]
	}
	return &mnode
}

func min(a int, b int) int {
	if a > b {
		return b
	}

	return a
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	nodes := []MerkleNode{} //存储所有默克尔树节点
	//构建叶子节点
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}
	j := 0
	// j代表的是某一层的第一个元素
	//第一层循环代表， nSize代表某一层的个数，每一循环一次减半
	for nSize := len(data); nSize > 1; nSize = (nSize + 1) / 2 {
		// 第二条循环 i+2 代表两两拼接。 i2是为了当个数是基数的时候，拷贝最后的元素。
		for i := 0; i < nSize; i += 2 {
			i2 := min(i+1, nSize-1) // 此处i2的值，只有当节点为奇数的时候，拷贝最后的元素i2=nSize-1，其他情况i2都是等于i+1
			node := NewMerkleNode(&nodes[j+i], &nodes[j+i2], nil)
			nodes = append(nodes, *node)
			//  j代表的是某一层的第一个元素
		}
		j += nSize
	}
	mTree := MerkleTree{&(nodes[len(nodes)-1])}
	return &mTree

}

/*
func main(){
	//mnode1:=MerkleNode{}
	//mnode2:=MerkleNode{}
	//
	//nodes:= []MerkleNode{}
	//nodes=append(nodes,mnode1)
	//nodes=append(nodes,mnode2)
	//fmt.Printf("%T\n",nodes)
	//字符串hash转换为字节
	hash1,_:=  hex.DecodeString("16f0eb42cb4d9c2374b2cb1de4008162c06fdd8f1c18357f0c849eb423672f5f")

	hash2,_:=  hex.DecodeString("cce2f95fc282b3f2bc956f61d6924f73d658a1fdbc71027dd40b06c15822e061")

	//大小端的转换
	tool.ReverseByte(hash1)  //大端模式显示转换为小端模式进行运算
	tool.ReverseByte(hash2)//大端模式显示转换为小端模式进行运算

	data:=[][]byte{
		hash1,
		hash2,
	}
	MerkleTree:=NewMerkleTree(data)//此处数据存储为小端模式
	tool.ReverseByte(MerkleTree.RootNode.Data)  //转换为大端模式，方便阅读
	fmt.Printf("%x",MerkleTree.RootNode.Data)
}*/
