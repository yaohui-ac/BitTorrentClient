bencode 协议格式
均为ascii码
数字 1234：i1234e
字符串 "hello"：5:hello # 长度:字符串内容
列表 l[数据1][数据2][数据3][…]e
    # []只是为了说明,实际上协议不用分隔符分开
字典 d[key1][value1][key2][value2][…]e
    # []只是为了说明,实际上协议不用分隔符分开,把key当成正常的bencode对象处理

字典实际上就是对应于某种结构体，key对应变量名
期望最终有一个类似于json.UnMarshal的功能,将字符串/解析成特定的结构体
特定结构体也可以Marshal成数据
