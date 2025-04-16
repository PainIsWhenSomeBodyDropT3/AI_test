-- SQL Tester Example Queries

-- 1. Calculate total sales for March 2024
SELECT SUM(amount) 
FROM orders 
WHERE order_date BETWEEN '2024-03-01' AND '2024-03-31';

-- 2. Find the top-spending customer
SELECT customer, SUM(amount) AS total_spent 
FROM orders 
GROUP BY customer 
ORDER BY total_spent DESC 
LIMIT 1;

-- 3. Calculate the average order value
SELECT AVG(amount) 
FROM orders;

-- 4. Group sales by month
SELECT 
    strftime('%Y-%m', order_date) AS month,
    SUM(amount) AS total_sales
FROM orders
GROUP BY month
ORDER BY month;

-- 5. Customer order statistics
SELECT 
    customer,
    COUNT(*) AS order_count,
    SUM(amount) AS total_spent,
    AVG(amount) AS avg_order_value,
    MIN(amount) AS min_order,
    MAX(amount) AS max_order
FROM orders
GROUP BY customer
ORDER BY total_spent DESC;

-- 6. Daily sales for March
SELECT 
    order_date,
    SUM(amount) AS daily_sales
FROM orders
WHERE order_date BETWEEN '2024-03-01' AND '2024-03-31'
GROUP BY order_date
ORDER BY order_date;

-- 7. Customers with orders above average
SELECT 
    customer,
    AVG(amount) as avg_order
FROM orders
GROUP BY customer
HAVING avg_order > (SELECT AVG(amount) FROM orders)
ORDER BY avg_order DESC;

-- 8. Running total of sales by date
SELECT 
    a.order_date,
    SUM(b.amount) as running_total
FROM orders a
JOIN orders b ON b.order_date <= a.order_date
GROUP BY a.order_date
ORDER BY a.order_date;

-- 9. Percentage of total sales by customer
SELECT 
    customer,
    SUM(amount) as total_spent,
    ROUND(SUM(amount) * 100.0 / (SELECT SUM(amount) FROM orders), 2) as percentage
FROM orders
GROUP BY customer
ORDER BY total_spent DESC;

-- 10. Orders with amount greater than customer's average
SELECT 
    o1.id,
    o1.customer,
    o1.amount,
    o1.order_date,
    (SELECT AVG(amount) FROM orders o2 WHERE o2.customer = o1.customer) AS customer_avg
FROM orders o1
WHERE o1.amount > (SELECT AVG(amount) FROM orders o2 WHERE o2.customer = o1.customer)
ORDER BY o1.customer, o1.order_date; 