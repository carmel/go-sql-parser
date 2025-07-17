package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {

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
