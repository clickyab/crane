package workers

const sql = "INSERT INTO ads ON DUPLICATE UPDATE SET a=5,b=12,c=VALUES(c)+423"
