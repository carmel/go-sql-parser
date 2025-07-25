package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseCreateTableWithColumnComments(t *testing.T) {
	sql := `CREATE TABLE users (
		id INT COMMENT 'User ID',
		name VARCHAR(255) COMMENT 'User Name',
		phone varchar(20) NOT NULL COMMENT 'Phone Number',
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

	if len(createTableStmt.TableSchema.Columns) != 3 {
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

	col3Def, ok := createTableStmt.TableSchema.Columns[2].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected second column to be a ColumnDef")
	}
	if col3Def.Comment == nil || col3Def.Comment.String() != "'Phone Number'" {
		t.Errorf("Expected comment 'Phone Number' for column name, but got %v", col3Def.Comment.String())
	}
	fmt.Println(createTableStmt.String())

}

func TestParseCreateTableWithTableComment(t *testing.T) {
	sql := `CREATE TABLE users (
		id INT
	) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT = 'User Table';`
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

	if createTableStmt.GetComment() != "User Table" {
		t.Errorf("Expected table comment 'User Table', but got %s", createTableStmt.GetComment())
	}

	fmt.Println(createTableStmt.String())

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

func TestParseColumnConstraints(t *testing.T) {
	sql := `CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(255) NOT NULL, email VARCHAR(255) NULL);`
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

	// id column
	idCol, ok := createTableStmt.TableSchema.Columns[0].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected id column to be a ColumnDef")
	}
	if !idCol.PrimaryKey {
		t.Errorf("Expected id column to be primary key")
	}

	// name column
	nameCol, ok := createTableStmt.TableSchema.Columns[1].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected name column to be a ColumnDef")
	}
	if nameCol.Nullable == nil {
		t.Errorf("Expected name column to be NULL")
	}

	// email column
	emailCol, ok := createTableStmt.TableSchema.Columns[2].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected email column to be a ColumnDef")
	}
	if emailCol.Nullable == nil {
		t.Errorf("Expected email column to be NULL")
	}
}

func TestParseColumnConstraintsExtended(t *testing.T) {
	sql := `CREATE TABLE products (
		id INT PRIMARY KEY AUTO_INCREMENT,
		sku VARCHAR(100) UNIQUE NOT NULL,
		price DECIMAL(10, 2) DEFAULT 0.00,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
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

	// id column
	idCol, ok := createTableStmt.TableSchema.Columns[0].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected id column to be a ColumnDef")
	}
	if !idCol.PrimaryKey {
		t.Errorf("Expected id column to be primary key")
	}
	if !idCol.AutoIncrement {
		t.Errorf("Expected id column to be auto_increment")
	}

	// sku column
	skuCol, ok := createTableStmt.TableSchema.Columns[1].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected sku column to be a ColumnDef")
	}
	if !skuCol.Unique {
		t.Errorf("Expected sku column to be unique")
	}
	if skuCol.Nullable == nil {
		t.Errorf("Expected sku column to be NULL")
	}

	// price column
	priceCol, ok := createTableStmt.TableSchema.Columns[2].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected price column to be a ColumnDef")
	}
	if priceCol.DefaultExpr == nil {
		t.Fatalf("Expected price column to have a default expression")
	}
	if priceCol.DefaultExpr.String() != "0.00" {
		t.Errorf("Expected price column default value to be '0.00', but got %s", priceCol.DefaultExpr.String())
	}

	// created_at column
	createdAtCol, ok := createTableStmt.TableSchema.Columns[3].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected created_at column to be a ColumnDef")
	}
	if createdAtCol.DefaultExpr == nil {
		t.Fatalf("Expected created_at column to have a default expression")
	}
	if createdAtCol.DefaultExpr.String() != "CURRENT_TIMESTAMP" {
		t.Errorf("Expected created_at column default value to be 'CURRENT_TIMESTAMP', but got %s", createdAtCol.DefaultExpr.String())
	}
	if createdAtCol.OnUpdate == nil {
		t.Fatalf("Expected created_at column to have ON UPDATE clause")
	}
	if createdAtCol.OnUpdate.String() != "CURRENT_TIMESTAMP" {
		t.Errorf("Expected created_at column ON UPDATE value to be 'CURRENT_TIMESTAMP', but got %s", createdAtCol.OnUpdate.String())
	}
}

func TestParseTableConstraints(t *testing.T) {
	sql := `CREATE TABLE order_items (
		order_id INT,
		product_id INT,
		quantity INT,
		PRIMARY KEY (order_id, product_id),
		UNIQUE KEY (product_id)
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

	// Primary Key
	pk, ok := createTableStmt.TableSchema.Columns[3].(*Key)
	if !ok {
		t.Fatalf("Expected primary key to be a Key")
	}
	if pk.Name != "PRIMARY KEY" {
		t.Errorf("Expected primary key name 'PRIMARY KEY', but got %s", pk.Name)
	}
	if len(pk.Columns.Items) != 2 {
		t.Fatalf("Expected primary key to have 2 columns")
	}
	if pk.Columns.Items[0].String() != "order_id" {
		t.Errorf("Expected first primary key column to be 'order_id', but got %s", pk.Columns.Items[0].String())
	}
	if pk.Columns.Items[1].String() != "product_id" {
		t.Errorf("Expected second primary key column to be 'product_id', but got %s", pk.Columns.Items[1].String())
	}

	// Unique Key
	uniqueKey, ok := createTableStmt.TableSchema.Columns[4].(*Key)
	if !ok {
		t.Fatalf("Expected unique key to be a Key")
	}
	if uniqueKey.Name != "UNIQUE KEY" {
		t.Errorf("Expected unique key name 'UNIQUE KEY', but got %s", uniqueKey.Name)
	}
	if len(uniqueKey.Columns.Items) != 1 {
		t.Fatalf("Expected unique key to have 1 column")
	}
	if uniqueKey.Columns.Items[0].String() != "product_id" {
		t.Errorf("Expected unique key column to be 'product_id', but got %s", uniqueKey.Columns.Items[0].String())
	}
}

func TestParseTableNameAndColumnNames(t *testing.T) {
	sql := `CREATE TABLE test_table (col1 INT, col2 VARCHAR(255));`
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

	// Table Name
	if createTableStmt.Identifier.Table.Name != "test_table" {
		t.Errorf("Expected table name 'test_table', but got %s", createTableStmt.Identifier.Table.Name)
	}

	// Column Names
	if len(createTableStmt.TableSchema.Columns) != 2 {
		t.Fatalf("Expected 2 columns, but got %d", len(createTableStmt.TableSchema.Columns))
	}

	col1, ok := createTableStmt.TableSchema.Columns[0].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected first column to be a ColumnDef")
	}
	if col1.Name.Ident.Name != "col1" {
		t.Errorf("Expected first column name 'col1', but got %s", col1.Name.Ident.Name)
	}

	col2, ok := createTableStmt.TableSchema.Columns[1].(*ColumnDef)
	if !ok {
		t.Fatalf("Expected second column to be a ColumnDef")
	}
	if col2.Name.Ident.Name != "col2" {
		t.Errorf("Expected second column name 'col2', but got %s", col2.Name.Ident.Name)
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
			if createTableStmt.Identifier != nil {
				fmt.Printf("Table Name: %s\n", createTableStmt.Identifier.Table.Name)
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
