from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
import time

email = "dianeburms1.6.1.990@gmail.com"  
password = "dianeburms1.6.1.990"         

driver = webdriver.Chrome()  

try:
    driver.get("https://www.semrush.com/login/")

    time.sleep(2)

    # Find email and password input fields and login button
    email_field = driver.find_element(By.NAME, "email")
    password_field = driver.find_element(By.NAME, "password")
    login_button = driver.find_element(By.XPATH, "//button[@type='submit']")

    # Enter credentials
    email_field.send_keys(email)
    password_field.send_keys(password)

    # Click the login button
    login_button.click()

    # Allow time for login to complete and page to load
    time.sleep(5)

    # Extract cookies
    cookies = driver.get_cookies()

    # Write cookies to cookieforce.txt
    with open("cookieforce.txt", "w") as f:
        f.write("# Netscape HTTP Cookie File\n")
        f.write("# http://curl.haxx.se/rfc/cookie_spec.html\n")
        for cookie in cookies:
            line = f"{cookie['domain']}\t{cookie['httpOnly']}\t{cookie['path']}\t{cookie['secure']}\t{int(cookie['expiry']) if 'expiry' in cookie else 0}\t{cookie['name']}\t{cookie['value']}\n"
            f.write(line)
    time.sleep(510)

    print("Cookies have been successfully saved to cookieforce.txt.")

except Exception as e:
    print("An error occurred:", e)

finally:
    time.sleep(5)
    # Close the driver
    driver.quit()
