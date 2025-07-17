package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseCreateTableWithColumnComments(t *testing.T) {
	sql := `CREATE TABLE users (
		id INT COMMENT 'User ID',
		name VARCHAR(255) COMMENT 'User Name'
	);`
	p := NewParser(sql)
	stmts, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse SQL: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, but got %d", len(stmts))
	}
	createTableStmt, ok := stmts[0].(*CreateTable)
	if !ok {
		t.Fatalf("Expected CreateTable statement, but got %T", stmts[0])
	}

	if len(createTableStmt.TableSchema.Columns) != 2 {
		t.Fatalf("Expected 2 columns, but got %d", len(createTableStmt.TableSchema.Columns))
	}

	col1Def, ok := createTableStmt.TableSchema.Columns[0].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected first column to be a ColumnDef")
	}
	if col1Def.Comment == nil || col1Def.Comment.String() != "'User ID'" {
		t.Errorf("Expected comment 'User ID' for column id, but got %v", col1Def.Comment.String())
	}

	col2Def, ok := createTableStmt.TableSchema.Columns[1].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected second column to be a ColumnDef")
	}
	if col2Def.Comment == nil || col2Def.Comment.String() != "'User Name'" {
		t.Errorf("Expected comment 'User Name' for column name, but got %v", col2Def.Comment.String())
	}
}

func TestParseCreateTableWithTableComment(t *testing.T) {
	sql := `CREATE TABLE users (
		id INT
	) COMMENT = 'User Table';`
	p := NewParser(sql)
	stmts, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse SQL: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, but got %d", len(stmts))
	}
	createTableStmt, ok := stmts[0].(*CreateTable)
	if !ok {
		t.Fatalf("Expected CreateTable statement, but got %T", stmts[0])
	}

	if len(createTableStmt.TableOptions) != 1 {
		t.Fatalf("Expected 1 table option, but got %d", len(createTableStmt.TableOptions))
	}

	tableOption := createTableStmt.TableOptions[0]
	if tableOption.Name.Name != "COMMENT" {
		t.Errorf("Expected table option name 'COMMENT', but got %s", tableOption.Name.Name)
	}

	comment, ok := tableOption.Value.(*StringLiteral)
	if !ok {
		t.Fatalf("Expected table option value to be a StringLiteral")
	}

	if comment.String() != "'User Table'" {
		t.Errorf("Expected table comment 'User Table', but got %s", comment.String())
	}
}

func TestParseCTEWithColumnAliases(t *testing.T) {
	sql := `WITH my_cte (col1, col2) AS (SELECT 1, 2) SELECT * FROM my_cte;`
	p := NewParser(sql)
	stmts, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse SQL: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, but got %d", len(stmts))
	}
	selectStmt, ok := stmts[0].(*SelectQuery)
	if !ok {
		t.Fatalf("Expected SelectQuery statement, but got %T", stmts[0])
	}

	if selectStmt.With == nil || len(selectStmt.With.CTEs) != 1 {
		t.Fatalf("Expected 1 CTE, but got %d", len(selectStmt.With.CTEs))
	}

	cte := selectStmt.With.CTEs[0]
	if len(cte.ColumnAliases) != 2 {
		t.Fatalf("Expected 2 column aliases, but got %d", len(cte.ColumnAliases))
	}

	if cte.ColumnAliases[0].Name != "col1" {
		t.Errorf("Expected column alias 'col1', but got %s", cte.ColumnAliases[0].Name)
	}

	if cte.ColumnAliases[1].Name != "col2" {
		t.Errorf("Expected column alias 'col2', but got %s", cte.ColumnAliases[1].Name)
	}
}

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
