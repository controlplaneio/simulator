import sqlite3
from sqlite3 import Error
import datetime
import argparse
from pathlib import Path

def create_connection(db_file):
    conn = None
    try:
        conn = sqlite3.connect(db_file)
    except Error as e:
        print(e)

    return conn


def check_orders(conn):
    current_date = datetime.date.today()

    cur = conn.cursor()
    cur.execute(f"SELECT date FROM orders where date = {current_date}")

    rows = cur.fetchall()

    if rows:
        print("Processing new order")
    else:
        print("No new orders to process")


def connect(db):
    # create a database connection
    conn = create_connection(db)
    with conn:
        check_orders(conn)

def main() -> None:
    parser = argparse.ArgumentParser(description='process new orders')
    parser.add_argument('--db',
                        metavar='db',
                        type=str,
                        default='/tmp/orders.db',
                        help='directory for orders database')
    args = parser.parse_args()
    location = Path(args.db)

    if location.is_file():
        try:
            connect(args.db)
        except KeyboardInterrupt:
            print("user interrupted the program.")
            raise SystemExit(0)
    else:
        print("Database does not exist")

if __name__ == '__main__':
    main()