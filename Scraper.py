import json
import os
import time
from dataclasses import dataclass
from concurrent.futures import ThreadPoolExecutor

import requests
from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.common import TimeoutException, NoSuchElementException
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import Select, WebDriverWait


@dataclass
class Card:
    """Class for storing information scraped from each card"""
    name: str
    cost: int
    power: int
    counter: int
    color: str
    card_type: str
    effect: str
    card_set: str
    attribute: str
    number: int
    img_path: str
    info: str

    def to_dict(self):
        """Convert card info to json format"""
        return {
            "name": self.name,
            "cost": self.cost,
            "power": self.power,
            "counter": self.counter,
            "color": self.color,
            "type": self.card_type,
            "effect": self.effect,
            "set": self.card_set,
            "attribute": self.attribute,
            "cardNo": self.number,
            "imgPath": self.img_path,
            "info": self.info
        }


class Scraper:
    def __init__(self):
        self.wait_time = 20.0
        # driver instance
        driver_options = webdriver.ChromeOptions()
        driver_options.add_argument("--headless")
        driver_options.add_argument("--no-sandbox")
        driver_options.add_argument("--disable-dev-shm-usage")
        self.driver = webdriver.Chrome(options=driver_options)
        self.driver.set_page_load_timeout(self.wait_time)

        # card information
        self.card_count = 0
        self.cards = []
        self.cards_json = []
        self.img_urls = []

    def driver_shutdown(self):
        self.driver.quit()

    def strip_header(self, info_element):
        """Remove header text and return attribute value"""
        info_element.h3.decompose()
        return info_element.get_text(strip=True)

    def extract_card_type(self, info_element):
        if info_element is None:
            return ""
        # Extract all the text from individual <span> elements
        spans = info_element.find_all("span")
        if len(spans) < 3:
            return ""
        # Extract the third <span> (card type) and process it
        card_type = spans[2].get_text(strip=True).upper()
        # Check if the card type is one of the valid types
        if card_type in ["CHARACTER", "LEADER", "STAGE", "EVENT"]:
            return card_type
        return ""

    def download_image(self, image_info):
        image_index, full_img_url = image_info
        img_data = requests.get(full_img_url).content
        image_filename = f"{image_index}.jpg"

        with open(image_filename, 'wb') as img_file:
            img_file.write(img_data)

    def scrape(self, url):
        # get url
        try:
            print(f"Getting {url}")
            self.driver.get(url)
            print(f"Got {url}")

            cookie_button = WebDriverWait(self.driver, self.wait_time).until(
                EC.presence_of_element_located((By.CLASS_NAME, "onetrust-close-btn-handler"))
            )
            cookie_button.click()
            print("Cookie button clicked!")
        except TimeoutException:
            print(f"Timeout while loading {url}")

        # start with dropdown visible
        options = self.driver.find_elements(By.XPATH, "//li[@class='selModalClose']")
        optionDropDown = WebDriverWait(self.driver, 20).until(
            EC.element_to_be_clickable((By.CLASS_NAME, "selModalButton"))
        )
        optionDropDown.click()

        for i in range(2, len(options) - 2):
            # dropdown select
            options = self.driver.find_elements(By.XPATH, "//li[@class='selModalClose']")
            if options[i].is_displayed() == True:
                options[i].click()
                time.sleep(2)
            # search this set
            search_button = WebDriverWait(self.driver, self.wait_time / 2).until(
                EC.element_to_be_clickable(
                    (By.XPATH, "//*[@id=\"cardlist\"]/main/article/div/div[1]/form/div[3]/input"))
            )
            search_button.click()
            # reopen dropdown menu for next iteration
            optionDropDown = WebDriverWait(self.driver, self.wait_time).until(
                EC.element_to_be_clickable((By.CLASS_NAME, "selModalButton"))
            )
            optionDropDown.click()

            page_source = self.driver.page_source
            soup = BeautifulSoup(page_source, "html.parser")
            results = soup.find(class_="resultCol")
            card_elements = results.find_all("dl", class_="modalCol")
            card_images = results.find_all("a", class_="modalOpen")

            set_name = soup.find(class_="selModalButton")
            print(f"--Viewing {set_name.get_text()}--")

            # scrape card info
            for card in card_elements:
                self.card_count += 1
                title_element = card.find("div", class_="cardName")
                cost_element = card.find("div", class_="cost")
                power_element = card.find("div", class_="power")
                counter_element = card.find("div", class_="counter")
                color_element = card.find("div", class_="color")
                type_element = card.find("div", class_="feature")
                effect_element = card.find("div", class_="text")
                set_element = card.find("div", class_="getInfo")
                attribute_element = card.find("div", class_="attribute")
                card_element = card.find("div", class_="infoCol")

                new_card = Card(
                    title_element.get_text(strip=True),
                    self.strip_header(cost_element),
                    self.strip_header(power_element),
                    self.strip_header(counter_element),
                    self.strip_header(color_element),
                    self.strip_header(type_element),
                    self.strip_header(effect_element),
                    self.strip_header(set_element),
                    self.strip_header(attribute_element),
                    self.card_count,
                    "assets/" + str(self.card_count) + ".jpg",
                    self.extract_card_type(card_element)
                )

                self.cards.append(new_card)
                self.cards_json.append(new_card.to_dict())

            # scrape images
            for idx, images in enumerate(card_images, start=len(self.img_urls) + 1):
                image = images.find("img", class_="lazy")
                image_link = image["data-src"]
                full_img_url = url.rstrip('/cardlist') + image_link.replace("..", "")
                self.img_urls.append((idx, full_img_url))

            print(f"Finished round with {self.card_count} total cards")

    def download_cards(self):
        """Create json of all card info and download all card images"""
        with open("data.json", 'w+', encoding='utf-8') as file:
            json.dump(self.cards_json, file, ensure_ascii=False, indent=4)
        print("JSON created")
        # move into card image directory
        os.chdir("./assets/cards/")
        num_workers = os.cpu_count() * 2
        print(f"Num of CPU: {num_workers}")
        with ThreadPoolExecutor(max_workers=num_workers) as executor:
            executor.map(self.download_image, self.img_urls)


scraper = Scraper()
start_time = time.time()
scraper.scrape("https://en.onepiece-cardgame.com/cardlist/")
scraper.driver_shutdown()
scraper.download_cards()
end_time = time.time()
print(f"Completed in: {end_time - start_time} seconds")
print(f"Scraped {scraper.card_count} cards")

