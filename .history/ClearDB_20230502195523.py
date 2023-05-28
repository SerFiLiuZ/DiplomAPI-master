import datetime
import mysql.connector
import time

while True:
    conn = mysql.connector.connect(
            host="localhost",
            user="root",
            password="1234",
            database="taskdb_v2"
        )

    cursor = conn.cursor()

    # Get the current time
    current_time = datetime.datetime.now()

    # Calculate the time difference threshold
    time_diff_threshold = datetime.timedelta(days=1)

    # Build the SQL query to delete records older than 24 hours
    sql = "DELETE FROM applications WHERE TIMESTAMPDIFF(SECOND, dataApplication, %s) > %s"
    cursor.execute(sql, (current_time, time_diff_threshold.total_seconds()))
    conn.commit()

    # Close the database connection
    cursor.close()
    conn.close()

    time.sleep(60*5)
    
    