import telebot
import json
import datetime
import mysql.connector

botToken= "6182612460:AAGQwtZa8TPijoa6YoiC5nkoWD3jhfpeRUI"
mode = ""
passwordOrg = '' 
fio = '' 
phoneNumber = '' 

mydb = mysql.connector.connect(
        host="localhost",
        user="root",
        password="1234",
        database="taskdb"
    )

bot = telebot.TeleBot(botToken)

@bot.message_handler(commands=['start'])
def start(message):
   global  mode 
   mode = "passwordOrg"
   mess = \
       f'Здравствуйте, <b>{message.from_user.first_name}</b>\n' \
       f'Чтобы оставить свою заявку, следуйте инструкциям:\n'\
       f'Введите пароль организации:\n'
 
   bot.send_message(message.chat.id, mess, parse_mode='html')

@bot.message_handler(content_types=['text'])
def get_text_messages(message):
    global mode 
    global passwordOrg  
    global fio
    global phoneNumber

    if mode == "passwordOrg":
        if checkPasswordOrg(message.text):
            bot.send_message(message.chat.id, f"Организация: {getNameOrg(message.text)}", parse_mode='html')
            bot.send_message(message.chat.id, f"Введите ФИО:", parse_mode='html')
            passwordOrg = message.text
            mode = "FIO"
        else:
            bot.send_message(message.chat.id, "Организации с таким паролем не существует. Проверьте правильность ввода пароля.", parse_mode='html')
    elif mode == "FIO":
        fio = message.text
        bot.send_message(message.chat.id, f"Введите номер телефона:", parse_mode='html')
        mode = "phoneNumder"
    elif mode == "phoneNumder":
        phoneNumber = message.text
        mess = \
            f'Введённные данные:\n' \
            f'ФИО: {fio}\n'\
            f'Телефонный номер: {phoneNumber}\n'
        bot.send_message(message.chat.id, mess, parse_mode='html')
        bot.send_message(message.chat.id, "Записано. Ваша заявка будет активна в течении 24 часов с этого момента.")
        bot.send_message(message.chat.id, "Если вы ввели не правильные данные, то попробуйте ещё раз, введя пароль организации.")
        mode = "passwordOrg"
        add_worker_in_db(passwordOrg, fio, phoneNumber, message.chat.id, message.from_user.username, datetime.datetime.now().replace(microsecond=0))

def checkPasswordOrg(password):
    # Create a cursor to execute SQL queries
    mycursor = mydb.cursor()

    # Execute the SQL query
    mycursor.execute(f"SELECT COUNT(*) FROM taskdb.users WHERE passwordCorp={password} and taskdb.users.status = 0")

    result = mycursor.fetchone()[0]

    # Commit the changes to the database
    mydb.commit()

    if result > 0:
        return True
    else:
        return False

def getNameOrg(password):
    # Create a cursor to execute SQL queries
    mycursor = mydb.cursor()

    # Execute the SQL query
    mycursor.execute(f"SELECT organizationName FROM users WHERE organizationPassword={password}")

    result = mycursor.fetchone()[0]

    # Commit the changes to the database
    mydb.commit()

    return result

def add_worker_in_db(passwordCorp, FIO, phoneNumber, chatID, userName, date):
    # Create a cursor to execute SQL queries
    mycursor = mydb.cursor()

    # Define the SQL query to insert the values into the table
    sql = "INSERT INTO applications (passwordCorp, FIO, phoneNumber, chatID, userName, dataApplication) VALUES (%s, %s, %s, %s, %s, %s)"
    values = (passwordCorp, FIO, phoneNumber, chatID, userName, date)

    # Execute the SQL query
    mycursor.execute(sql, values)

    # Commit the changes to the database
    mydb.commit()

    # Print a message to confirm the insertion
    print(mycursor.rowcount, "record inserted.")
    return

bot.polling()