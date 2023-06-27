# 脚本说明

对指定词条目录下，对冲突词条（key值相同、value不同）、有效词条 （key相同且value值相同、或者 key值不同的词条）以及 key相同且value值相同词条进行筛选，并打印出词条数量，便于统计验证。

# 使用说明

```shell
./main --filePath 指定词条路径  --region 指定语言目录 --outPath 指定输出目录(会生成conflict.json[冲突词条]、efficient.json[有效词条]以及keyEqualvalue.json[key和value都相同词条]文件)
参数解释：
filePath 指定词条目录所在路径
region 指定对目录下哪种语言的词条包进行查找
outPath 输出目录
```

