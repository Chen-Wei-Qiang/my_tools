# 脚本说明

对project-api、wiki-api、project-web以及wiki-web中，对冲突词条（key值相同、value不同）和 有效词条 （key相同且value值相同、或者 key值不同的词条）进行筛选。

# 目录结构

```shell
mytools
- citiao
-- project-api // project-api词条目录下有中英日词条
--- en.json
--- ja.json
--- zh.json
-- project-web // project-web词条目录下有中英日词条
--- en.json
--- ja.json
--- zh.json
-- wiki-api    // wiki-api词条目录下有中英日词条
--- en.json
--- ja.json
--- zh.json
-- wiki-web    // wiki-web词条目录下有中英日词条
--- en.json
--- ja.json
--- zh.json
- output
-- conflict.json  // 冲突词条
-- efficient.json // 有效词条
- main.go
```

# 使用说明

```
./main --filePath 指定词条路径  --region 指定语言目录
参数解释：
filePath 指定目录结构中 citiao目录 路径
region 指定是对project-api、wiki-api、project-web以及wiki-web 下 哪种语言来进行查找
```

