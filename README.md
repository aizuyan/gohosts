# gohosts
terminal hosts manager


命令行hosts管理工具

## 使用
下载`dist`中对应的版本，复制到/usr/bin/gohost（或其他PATH路径下）
```shell
sudo gohosts
```

上面的命令进入软件。
首次进入会进行初始化，已有的hosts会保存到一个叫`backup`的hosts item中，默认打开。

### 命令
```
1. [tab] => 切换左右view
2. [shift + a] => 聚焦左view时候，添加新的hosts item，新添加的hosts item，默认关闭
3. [Enter] => 聚焦输入框view时候，回车，终止输入并新建hosts item，如果有错误，文字呈红色；无内容时候，回车，移出输入框view
4. [↑] => 聚焦左view时候，焦点向上移动一行
5. [↓] => 聚焦左view时候，焦点向下移动一行
6. [←] => 聚焦左view的时候，关闭一个hosts item
7. [→] => 聚焦左view的时候，打开一个hosts item
```

## mac 效果
![](./gohosts_32.gif)

![](./images/mac.png)
