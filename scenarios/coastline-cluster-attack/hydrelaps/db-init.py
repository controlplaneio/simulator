import sqlite3

from faker import Faker

def main():
    conn = sqlite3.connect('/host/users.db')
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
        sqlite_create_table = """CREATE TABLE USERS(
                                ID           INT PRIMARY KEY     NOT NULL,
                                IP           TEXT    NOT NULL,
                                HOSTNAME     INT     NOT NULL,
                                MAC_ADDRESS  TEXT    NOT NULL,
                                USER_AGENT   BLOB    NOT NULL);"""

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

        for _ in range (30):
            id = _
            ip = fake.ipv4_public()
            hostname = fake.hostname()
            mac_address = fake.mac_address()
            user_agent = fake.user_agent()

            data = str(f"{id}, '{ip}', '{hostname}', '{mac_address}', '{user_agent}'")
            print(data)

            sqlite_insert_query = f"""INSERT INTO USERS
                                (id, ip, hostname, mac_address, user_agent)
                                VALUES
                                ({data})"""
            print(sqlite_insert_query)

            count = cursor.execute(sqlite_insert_query)
            conn.commit()
            print("Record inserted successfully into SqliteDb_developers table ", cursor.rowcount)

    except sqlite3.Error as error:
        print("Failed to insert data into sqlite table", error)

if __name__ == "__main__":
    main()

