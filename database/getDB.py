import yfinance as yf
import requests
import pandas as pd
from datetime import datetime

YEARS_PRE_BINANCE = 11  # Años a obtener de Yahoo Finance (2013-2024)


def get_yahoo_data():
    try:
        start_date = "2013-01-01"
        end_date = f"{2013 + YEARS_PRE_BINANCE}-12-31"
        print(f"Obteniendo datos de Yahoo Finance ({
              start_date} a {end_date})...")
        btc = yf.download(
            "BTC-USD",
            start=start_date,
            end=end_date,
            progress=False,
            timeout=30
        )
        if btc.empty:
            raise ValueError("No se encontraron datos en Yahoo Finance")
        btc.reset_index(inplace=True)
        btc["fecha_apertura"] = btc["Date"].dt.strftime("%d/%m/%Y 00:00")
        btc["fecha_cierre"] = btc["Date"].dt.strftime("%d/%m/%Y 23:59")
        if isinstance(btc.columns, pd.MultiIndex):
            btc.columns = ['_'.join(col).strip() if col[1] else col[0]
                           for col in btc.columns.values]
        return btc[["fecha_apertura", "Open_BTC-USD",
                    "fecha_cierre", "Close_BTC-USD"]].rename(
            columns={"Open_BTC-USD": "open", "Close_BTC-USD": "close"}
        )
    except Exception as e:
        print(f"Error con Yahoo Finance: {str(e)}")
        return pd.DataFrame()


# --- 2. Obtener datos de Binance (2025-actualidad) ---
def get_binance_data():
    try:
        print("Obteniendo datos de Binance (2025-actualidad)...")
        start_date = datetime(2025, 1, 1)
        start_timestamp = int(start_date.timestamp() * 1000)
        url = f"https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=1d&startTime={
            start_timestamp}&limit=1000"
        response = requests.get(url, timeout=10)
        response.raise_for_status()
        data = response.json()
        df = pd.DataFrame(data, columns=[
            'timestamp', 'open', 'high', 'low', 'close', 'volume',
            'close_time', 'quote_volume', 'count', 'taker_buy_volume',
            'taker_buy_quote_volume', 'ignore'
        ])
        df["timestamp"] = pd.to_datetime(df["timestamp"], unit="ms")
        df["fecha_apertura"] = df["timestamp"].dt.strftime("%d/%m/%Y 00:00")
        df["fecha_cierre"] = df["timestamp"].dt.strftime("%d/%m/%Y 23:59")
        return df[["fecha_apertura", "open", "fecha_cierre", "close"]].astype({
            "open": float,
            "close": float
        })
    except Exception as e:
        print(f"Error con Binance API: {str(e)}")
        return pd.DataFrame()


# --- 3. Procesar y combinar datos ---
def process_data():
    df_yahoo = get_yahoo_data()
    df_binance = get_binance_data()
    if df_yahoo.empty or df_binance.empty:
        raise ValueError("No se pudieron obtener todos los datos necesarios")
    df_full = pd.concat([df_yahoo, df_binance], ignore_index=True)
    print("\nColumnas finales en df_full:", df_full.columns.tolist())
    df_full = df_full.dropna(axis=1, how='all')
    df_full["fecha_datetime"] = pd.to_datetime(
        df_full["fecha_apertura"], dayfirst=True
    )
    df_full = df_full.sort_values(
        "fecha_datetime").drop("fecha_datetime", axis=1)
    df_apertura = df_full[["fecha_apertura", "open"]].rename(
        columns={"fecha_apertura": "fecha", "open": "apertura"}
    )
    df_cierre = df_full[["fecha_cierre", "close"]].rename(
        columns={"fecha_cierre": "fecha", "close": "cierre"}
    )
    return df_apertura, df_cierre


# --- 4. Exportar a CSV ---
def export_to_csv(df_apertura, df_cierre):
    try:
        file_apertura = "bitcoin_apertura_2013_actual.csv"
        file_cierre = "bitcoin_cierre_2013_actual.csv"
        df_apertura.to_csv(file_apertura, index=False,
                           sep=",", encoding="utf-8")
        df_cierre.to_csv(file_cierre, index=False, sep=",", encoding="utf-8")
        print("Archivos exportados con éxito:")
        print(f"- {file_apertura}")
        print(f"- {file_cierre}")
        print("Total de registros: ", len(
            df_apertura), "dias (2013-actualidad)")
    except Exception as e:
        print(f"Error al exportar CSV: {str(e)}")


# --- Ejecución principal ---
if __name__ == "__main__":
    print("\n🔍 Iniciando descarga de datos históricos de Bitcoin (2013-actualidad)...")
    try:
        df_apertura, df_cierre = process_data()
        export_to_csv(df_apertura, df_cierre)
    except Exception as e:
        print(f"\n❌ Error crítico: {str(e)}")
