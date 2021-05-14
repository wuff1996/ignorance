### The Flow Of Go-Ethereum Source Code

----------



#### 需要的工具:

- Goland  Software（查看源码） 

- Go-ethereum Source Code: 

  https://github.com/r1cs/go-ethereum

- Other References:

  https://github.com/r1cs/go-ethereum-code-analysis

  

  - **Email**: 1027309561@qq.com

  ##### 本文风格： 

  1：中英文以空格符分开

  2：默认无结尾符号

  3：加粗代表关键词

  

  ##### 目标受众：

  热爱区块链技术，已经对区块链技术有系统性的解且想从事或研究底层逻辑的人。

  

  ----------------------------

  

  #### 内容大概

  - 技术大概(Technical Protocol)
  - 内容逻辑(Over Flow of Source Code)
  - ETH 层细节(Eth Protocol Layer)
  - P2P 层细节(P2P Protocol Layer)
  - 一些有趣的 tips(functions & modules)

  

  ----------------

  #### 技术大概(Technical Protocol)

  

  在**以太坊白皮书**中(https://gitbub.com/r1cs/ignorance/BlockChain/book/eth_whitepaper.pdf)粗略介绍了以下几种概念

  

  1: 账户(Account)

  2: 消息与交易(Msg & Transaction)

  3: 状态转移(State Transition)

  4: 智能合约(Smart Contract)

  5: 区块链与挖矿(Blockchain & Mining)

  

  -------------

  1: 账户(Account)

  

  在讲账户之前，需要讲一下**Merkle Partricia Tree**，所有的账户都被保存在此结构(Key Value)中，这种**数据结构**保证了以太坊的**高效**搜索能力且能**快速证明**该**状态树**是否有效。

  

  由于网络上可以搜到很多这种资料，因此就不再赘述

  ```txt
  参考链接: https://ethfans.org/toya/articles/588
  ```

  

  所以账户就是以 地址(address)为key，

  ```go
  //location: crypto/crypto.go:276 line
  
  func PubkeyToAddress(p ecdsa.PublicKey) common.Address {
  	pubBytes := FromECDSAPub(&p)
  	return common.BytesToAddress(Keccak256(pubBytes[1:])[12:]) //取最后的20byte
  }
  ```

  

  保存着以下Value: [**Balance, Nonce, Contract_Root**]的**RLP**值，Balance代表账户余额；Nonce代表该地址的交易数，使得每次交易的**数字签名(ECDSA)**唯一；Contract_Root为帕特里夏树根，在该地址为一个合约控制的情况下包含该合约的内存。

  

  其中新出现了一些名词。**RLP**，递归长度前缀是一中编码算法，用于编码任意嵌套结构的对象，它是以太坊中数据序列化的主要方法。

  ```txt
  以太坊中的所有数据都以递归长度前缀编码(Recursive Length Prefix Encoding)形式存储
  其基本的思想是把数据类型和长度编码成一个单独的字节放在实际数据的前面（ 例如‘dog’ 的字节数组编码为[ 100, 111, 103 ], 于
  是串接后就成了[ 67, 100, 111, 103 ].） 注意 RLP 编码正如其名字表示的一样， 是递归的； 当 RLP 编码一个数组时， 实际上是在对每
  一个元素的 RLP 编码级联成的字符串编码.。 需要进一步提请注意的是， 以太坊中所有数据都是整数
  
  参考资料： https://github.com/r1cs/go-ethereum-code-analysis/blob/master/rlp-analysis.md
  ```

  ```go
  //location: rlp/encode.go:53 line
  
  // Encode writes the RLP encoding of val to w. Note that Encode may
  // perform many small writes in some cases. Consider making w
  // buffered.
  //
  // Please see package-level documentation of encoding rules.
  func Encode(w io.Writer, val interface{}) error {
  	if outer, ok := w.(*encbuf); ok {
  		// Encode was called by some type's EncodeRLP.
  		// Avoid copying by writing to the outer encbuf directly.
  		return outer.encode(val) //由于go语言没有方法重载和通用类型，所以此方法通过反射拿到底层类型的encoder函数进行编码
  	}
  	eb := encbufPool.Get().(*encbuf)
  	defer encbufPool.Put(eb)
  	eb.reset()
  	if err := eb.encode(val); err != nil {
  		return err
  	}
  	return eb.toWriter(w)
  }
  
  ```

  ，我刚才得知已经有人写过了，所以我就不再赘述了技术大概了，参考资料：

  ```txt
  https://github.com/ZtesoftCS/go-ethereum-code-analysis
  ```

  

  

  ---------------------

  #### 内容逻辑(Over Flow of Source Code)

  entry point

  ```go
  //location: cmd/geth/main.go:332 line
  
  // geth is the main entry point into the system if no special subcommand is ran.
  // It creates a default node based on the command line arguments and runs it in
  // blocking mode, waiting for it to be shut down.
  func geth(ctx *cli.Context) error {
  	if args := ctx.Args(); len(args) > 0 {
  		return fmt.Errorf("invalid command: %q", args[0])
  	}
  
      prepare(ctx)	//操作内存缓存分配上限，启动metric系统(用于监测各个操作的资源消耗)
  	stack, backend := makeFullNode(ctx)	//加载配置文件，注册以太坊服务，API接口
  	defer stack.Close()
  
  	startNode(ctx, stack, backend)	//启动节点和所有以注册的协议，然后解锁相关账户并开启RPC/IPC接口，最后根据配置文件启动挖矿。
  	stack.Wait()	//<-n.stop 阻塞监听关闭通道
  	return nil
  }
  
  ```

  

  

  

  ---------------

  #### ETH 层细节

  这里主要讲downloader，fetcher，转发交易。

  downloader:

  

  

  

  -----------------

  #### P2P 层细节

  主要讲节点发现，信息传输。

  

  