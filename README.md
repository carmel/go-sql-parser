# Golang SQL Parser Library / Golang SQL 解析库

[English](#english) | [中文](#中文)

---

## English

### Overview

This is a Golang-based SQL parsing library inspired by the open-source [ClickHouse SQL Parser](https://github.com/AfterShip/clickhouse-sql-parser) with significant improvements and extensions. The library provides comprehensive support for modern SQL syntax including AUTO_INCREMENT, PRIMARY KEY, INDEX, UNIQUE, COMMENT, DEFAULT, ON UPDATE, and many other features.

### Features

- **🚀 High Performance**: Optimized parsing engine for fast SQL statement processing
- **🔧 Extended SQL Support**: Enhanced support for MySQL-style syntax extensions
- **📝 Comprehensive Coverage**: Handles complex SQL statements with advanced features
- **🎯 Easy Integration**: Simple API design for seamless integration into Go projects
- **🔍 Detailed AST**: Provides detailed Abstract Syntax Tree for advanced analysis

#### Supported SQL Features

- **Table Operations**: CREATE TABLE, DROP TABLE, ALTER TABLE
- **Index Management**: CREATE INDEX, DROP INDEX, UNIQUE INDEX
- **Column Constraints**:
  - PRIMARY KEY
  - UNIQUE constraints
  - DEFAULT values
  - AUTO_INCREMENT
  - ON UPDATE triggers
  - COMMENT annotations
- **Data Types**: All standard SQL data types with MySQL extensions
- **Query Operations**: SELECT, INSERT, UPDATE, DELETE with complex conditions
- **Advanced Features**: Subqueries, JOINs, CTEs, window functions

### Installation

```bash
go get github.com/carmel/go-sql-parser
```

### Quick Start

```go
package main

import (
    "fmt"
    "github.com/carmel/go-sql-parser/parser"
)

func main() {

	// 示例SQL脚本
	query := `
-- 用户信息表
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    email VARCHAR(100) UNIQUE COMMENT '邮箱地址',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    phone VARCHAR(20) COMMENT '手机号',
    status TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
insert into users(username,phone)values('Jhon','131');
/*
订单详情表
*/
CREATE TABLE orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '订单ID',
    user_id INT NOT NULL COMMENT '用户ID',
    order_no VARCHAR(32) NOT NULL COMMENT '订单号',
    total_amount DECIMAL(10,2) NOT NULL COMMENT '总金额',
    status TINYINT DEFAULT 0 COMMENT '订单状态：0-待支付，1-已支付，2-已完成，3-已取消',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';
`

	parser := NewParser(query)
	// Parse query into AST
	statements, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}
	for _, stmt := range statements {
		if createTableStmt, ok := stmt.(*CreateTable); ok {
			// Extract table name
			if createTableStmt.Name != nil {
				fmt.Printf("Table Name: %s\n", createTableStmt.Name)
			}

			// Extract table comment
			for _, option := range createTableStmt.TableOptions {

				if option.Name != nil && strings.ToUpper(option.Name.Name) == "COMMENT" {
					if option.Value != nil {
						if strLiteral, ok := option.Value.(*StringLiteral); ok {
							fmt.Printf("Table Comment: %s\n", strLiteral.Literal)
						}
					}
				}
			}

			fmt.Println("\nColumns:")
			// Extract column names and comments
			if createTableStmt.TableSchema != nil && createTableStmt.TableSchema.Columns != nil {
				for _, column := range createTableStmt.TableSchema.Columns {
					fmt.Println(column)
				}
			}
		}
	}
}
```

### API Reference

#### Core Functions

```go
// Parse a single SQL or mutiple statement
func Parse() ([]Expr, error)
```

#### Statement Types

- `CreateTableStmt`: CREATE TABLE statements
- `DropTableStmt`: DROP TABLE statements
- `AlterTableStmt`: ALTER TABLE statements
- `SelectStmt`: SELECT queries
- `InsertStmt`: INSERT statements
- `UpdateStmt`: UPDATE statements
- `DeleteStmt`: DELETE statements

### Advanced Usage

#### AST Traversal

```go
// Implement visitor pattern for AST traversal
type MyVisitor struct{}

func (v *MyVisitor) VisitCreateTable(stmt *parser.CreateTableStmt) {
    // Handle CREATE TABLE statements
}

func (v *MyVisitor) VisitSelectStmt(stmt *parser.SelectStmt) {
    // Handle SELECT statements
}

// Use visitor
visitor := &MyVisitor{}
stmt.Accept(visitor)
```

### Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Acknowledgments

- Inspired by [AfterShip/clickhouse-sql-parser](https://github.com/AfterShip/clickhouse-sql-parser)
- Thanks to all contributors and the open-source community

---

## 中文

### 概述

这是一个基于 Golang 开发的 SQL 解析库，受到开源项目 [ClickHouse SQL Parser](https://github.com/AfterShip/clickhouse-sql-parser) 的启发，并进行了重大改进和扩展。该库全面支持现代 SQL 语法，包括 AUTO_INCREMENT、PRIMARY KEY、INDEX、UNIQUE、COMMENT、DEFAULT、ON UPDATE 等多种特性。

### 特性

- **🚀 高性能**: 优化的解析引擎，快速处理 SQL 语句
- **🔧 扩展的 SQL 支持**: 增强对 MySQL 风格语法扩展的支持
- **📝 全面覆盖**: 处理具有高级特性的复杂 SQL 语句
- **🎯 易于集成**: 简单的 API 设计，无缝集成到 Go 项目中
- **🔍 详细的 AST**: 提供详细的抽象语法树用于高级分析

#### 支持的 SQL 特性

- **表操作**: CREATE TABLE、DROP TABLE、ALTER TABLE
- **索引管理**: CREATE INDEX、DROP INDEX、UNIQUE INDEX
- **列约束**:
  - PRIMARY KEY（主键）
  - UNIQUE 约束（唯一约束）
  - DEFAULT 值（默认值）
  - AUTO_INCREMENT（自增）
  - ON UPDATE 触发器
  - COMMENT 注释
- **数据类型**: 所有标准 SQL 数据类型及 MySQL 扩展
- **查询操作**: SELECT、INSERT、UPDATE、DELETE 及复杂条件
- **高级特性**: 子查询、JOIN、CTE、窗口函数

### 安装

```bash
go get github.com/carmel/go-sql-parser
```

### 快速开始

```go
package main

import (
    "fmt"
    "github.com/carmel/go-sql-parser/parser"
)

func main() {
// 示例SQL脚本
	query := `
-- 用户信息表
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    email VARCHAR(100) UNIQUE COMMENT '邮箱地址',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    phone VARCHAR(20) COMMENT '手机号',
    status TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
insert into users(username,phone)values('Jhon','131');
/*
订单详情表
*/
CREATE TABLE orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '订单ID',
    user_id INT NOT NULL COMMENT '用户ID',
    order_no VARCHAR(32) NOT NULL COMMENT '订单号',
    total_amount DECIMAL(10,2) NOT NULL COMMENT '总金额',
    status TINYINT DEFAULT 0 COMMENT '订单状态：0-待支付，1-已支付，2-已完成，3-已取消',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';
`

	parser := NewParser(query)
	// Parse query into AST
	statements, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}
	for _, stmt := range statements {
		if createTableStmt, ok := stmt.(*CreateTable); ok {
			// Extract table name
			if createTableStmt.Name != nil {
				fmt.Printf("Table Name: %s\n", createTableStmt.Name)
			}

			// Extract table comment
			for _, option := range createTableStmt.TableOptions {

				if option.Name != nil && strings.ToUpper(option.Name.Name) == "COMMENT" {
					if option.Value != nil {
						if strLiteral, ok := option.Value.(*StringLiteral); ok {
							fmt.Printf("Table Comment: %s\n", strLiteral.Literal)
						}
					}
				}
			}

			fmt.Println("\nColumns:")
			// Extract column names and comments
			if createTableStmt.TableSchema != nil && createTableStmt.TableSchema.Columns != nil {
				for _, column := range createTableStmt.TableSchema.Columns {
					fmt.Println(column)
				}
			}
		}
	}
}
```

#### 核心函数

```go
// 解析单个或多个 SQL 语句
func Parse() ([]Expr, error)
```

#### 语句类型

- `CreateTableStmt`: CREATE TABLE 语句
- `DropTableStmt`: DROP TABLE 语句
- `AlterTableStmt`: ALTER TABLE 语句
- `SelectStmt`: SELECT 查询
- `InsertStmt`: INSERT 语句
- `UpdateStmt`: UPDATE 语句
- `DeleteStmt`: DELETE 语句

### 高级用法

#### AST 遍历

```go
// 实现访问者模式进行 AST 遍历
type MyVisitor struct{}

func (v *MyVisitor) VisitCreateTable(stmt *parser.CreateTableStmt) {
    // 处理 CREATE TABLE 语句
}

func (v *MyVisitor) VisitSelectStmt(stmt *parser.SelectStmt) {
    // 处理 SELECT 语句
}

// 使用访问者
visitor := &MyVisitor{}
stmt.Accept(visitor)
```

### 许可证

本项目采用 MIT 许可证 - 详情请查看 [LICENSE](LICENSE) 文件。

### 致谢

- 受到 [AfterShip/clickhouse-sql-parser](https://github.com/AfterShip/clickhouse-sql-parser) 的启发
- 感谢所有贡献者和开源社区
