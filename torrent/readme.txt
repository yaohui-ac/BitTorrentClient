根据torrent文件不同,存在单文件模式和多文件模式

为了对外接口统一, 统一转化为torrentFile结构
后续网络通信 只需要面向torrentFile即可
读取文件 ->计算哈希 -> 填充结构

tracker通信 http/udp通信

peer 通信
1.握手 交换信息
    part1:  1byte: 第二块长度
    part2:  长度可变: 协议信息(BitTorrent Protocol)
    part3:  8bytes 保留位(全部为0)
    part4：  20bytes infoSHA
    part5:  20bytes peers id

2.获取peer数据情况
3.指定piece下载