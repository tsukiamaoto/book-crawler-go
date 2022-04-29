# Book-Crawler-Go

> 用來爬取書店的資料。現階段只有抓取博客來。

### Feature
* 自動抓取博客來網站
* 排定每小時執行一次

### Architecture
![結構圖](https://github.com/tsukiamaoto/book-crawler-go/blob/master/book%20store%20architeture.png?raw=true)

### Installation And Usage

``` bash
# Download
git clone https://github.com/tsukiamaoto/book-crawler

# Install go package
go mod download

# Run
go run main.go

```

### TODO

- [x] 掛載proxy
- [ ] Docker 部署
- [ ] 抓更多種類的書籍

### Problem

1. 博客來的網站會阻擋國外IP訪問國內頁面

2. ~~同一個IP不能訪問超過20-50次，會被鎖IP一分鐘~~
> 掛載proxy解決
3. ~~國外IP從首頁進來，會連到海外專區~~
> 不從首頁訪問