import streamlit as st
import requests

def main():
    st.title('Order management')

    input_text = st.text_input('Введите uid заказа:')

    if st.button('Найти заказ'):
        url = 'http://localhost:8080/orders'

        response = requests.post(url, json={'id': input_text})

        if response.status_code == 200:
            result = response.json()
            st.write('Информация по заказу:', result)
        else:
            st.error(f'Ошибка выполнения POST запроса: {response.status_code}')

if __name__ == "__main__":
    main()
