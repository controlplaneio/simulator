import sqlite3

from faker import Faker

def main():
    conn = sqlite3.connect('/tmp/orders.db')
    cursor = conn.cursor()
    print("Successfully Connected to SQLite")

    create_table(cursor, conn)
    insert_data(cursor, conn)

    cursor.close()
    if conn:
            conn.close()
            print("The SQLite connection is closed")

def create_table(cursor, conn):
    try:
        cursor = cursor
        conn = conn
        sqlite_create_table = '''CREATE TABLE orders (
                                id INT PRIMARY KEY,
                                fullname TEXT NOT NULL,
                                address BLOB NOT NULL,
                                email TEXT NOT NULL,
                                date TEXT NOT NULL,
                                credit_card INT NOT NULL,
                                cc_expire TEXT NOT NULL,
                                orders TEXT)'''

        cursor.execute(sqlite_create_table)
        conn.commit()
        print("SQLite table created")

    except sqlite3.Error as error:
        print("Failed to insert data into sqlite table", error)

def insert_data(cursor, conn):
    try:
        cursor = cursor
        conn = conn
        fake = Faker('en_US')

        items = ["frying pan", "mobile phone charger", "usb cable", "salt", "laptop charger", "eggs", "milk", "bread", "printer toner", "usb key", "gift card", "birthday card", "tooth brush", "flag_ctf{SEC_INCIDENT_CUSTOMER_DATA_DISCOVERED}","usb cable",  "laptop charger", "shower gel", "gift card"]

        for id, item in enumerate(items):

            id = id
            first_name = fake.first_name()
            first_ab = (first_name[0].lower())
            last_name = fake.last_name()
            last_lower = (last_name.lower())
            full_name = first_name+" "+last_name
            domain = fake.free_email_domain()
            login = first_ab+last_lower
            email = login+'@'+domain
            address = fake.address()
            date = fake.date_between(start_date='-1y',end_date='today')
            credit_card = fake.credit_card_number()
            cc_expire = fake.credit_card_expire()
            orders = item

            data = str(f"{id}, '{full_name}', '{address}', '{email}', '{date}', {credit_card}, '{cc_expire}', '{orders}'")
            print(data)

            sqlite_insert_query = f"""INSERT INTO ORDERS
                                (id, fullname, address, email, date, credit_card, cc_expire, orders)
                                VALUES
                                ({data})"""
            print(sqlite_insert_query)

            count = cursor.execute(sqlite_insert_query)
            conn.commit()
            print("Record inserted successfully into sqlite orders table ", cursor.rowcount)

    except sqlite3.Error as error:
        print("Failed to insert data into sqlite table", error)

if __name__ == "__main__":
    main()