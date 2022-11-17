根据torrent文件不同,存在单文件模式和多文件模式

为了对外接口统一, 统一转化为torrentFile结构
后续网络通信 只需要面向torrentFile即可
读取文件 ->计算哈希 -> 填充结构
