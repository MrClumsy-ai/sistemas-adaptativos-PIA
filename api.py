from flask import Flask, request, jsonify
import pandas as pd
import numpy as np
from keras.models import load_model
import joblib

app = Flask(__name__)

# Cargar modelo y recursos al iniciar la aplicación
scaler_clausura = joblib.load('./modelos/scaler_clausura.save')
model_clausura = load_model('./modelos/mejor_modelo_clausura.h5')
scaler_apertura = joblib.load('./modelos/scaler_apertura.save')
model_apertura = load_model('./modelos/mejor_modelo_apertura.h5')
ventana = 8  # Tamaño de ventana temporal


def generar_ventanas(datos_preproc, ventana):
    """Genera ventanas temporales para el modelo"""
    X = []
    for i in range(len(datos_preproc) - ventana + 1):
        X.append(datos_preproc[i:i+ventana, 0])
    return np.array(X)


@app.route('/predict_apertura', methods=['POST'])
def predict_apertura():
    """
    Endpoint para hacer predicciones
    Espera un JSON con:
    {
        "data": [valores_numericos],
        "dates": [fechas_opcional]
    }
    """
    try:
        # Obtener datos del request
        request_data = request.get_json()
        input_data = np.array(request_data['data']).reshape(-1, 1)
        # Preprocesamiento
        datos_estandarizados = scaler_apertura.transform(input_data)
        # Generar ventanas temporales
        X = generar_ventanas(datos_estandarizados, ventana)
        X = X.reshape(X.shape[0], X.shape[1], 1)
        # Hacer predicciones
        y_pred = model_apertura.predict(X)
        y_pred_inverse = scaler_apertura.inverse_transform(y_pred)
        # Preparar respuesta
        response = {
            "predictions": y_pred_inverse.flatten().tolist()
        }
        # Agregar fechas si se proporcionaron
        if 'dates' in request_data and len(request_data[
                'dates']) == len(input_data):
            prediction_dates = request_data['dates'][ventana:]
            response["dates"] = prediction_dates
        return jsonify(response)
    except Exception as e:
        return jsonify({"error": str(e)}), 400


@app.route('/predict_clausura', methods=['POST'])
def predict_clausura():
    """
    Endpoint para hacer predicciones
    Espera un JSON con:
    {
        "data": [valores_numericos],
        "dates": [fechas_opcional]
    }
    """
    try:
        # Obtener datos del request
        request_data = request.get_json()
        input_data = np.array(request_data['data']).reshape(-1, 1)
        # Preprocesamiento
        datos_estandarizados = scaler_clausura.transform(input_data)
        # Generar ventanas temporales
        X = generar_ventanas(datos_estandarizados, ventana)
        X = X.reshape(X.shape[0], X.shape[1], 1)
        # Hacer predicciones
        y_pred = model_clausura.predict(X)
        y_pred_inverse = scaler_clausura.inverse_transform(y_pred)
        # Preparar respuesta
        response = {
            "predictions": y_pred_inverse.flatten().tolist()
        }
        # Agregar fechas si se proporcionaron
        if 'dates' in request_data and len(request_data[
                'dates']) == len(input_data):
            prediction_dates = request_data['dates'][ventana:]
            response["dates"] = prediction_dates
        return jsonify(response)
    except Exception as e:
        return jsonify({"error": str(e)}), 400


@app.route('/predict_last_apertura', methods=['GET'])
def predict_last_apertura():
    """
    Endpoint para predecir el próximo valor basado en los últimos datos del CSV
    """
    try:
        # Cargar datos históricos
        data_file = pd.read_csv('./database/corr_bitcoin_diario_apertura.csv')
        datos = data_file['apertura'].to_numpy().reshape(-1, 1)
        # Tomar los últimos 'ventana' valores
        last_window = datos[-ventana:]
        last_window_scaled = scaler_apertura.transform(last_window)
        # Preparar para predicción
        X = last_window_scaled.T.reshape(1, ventana, 1)
        # Hacer predicción
        y_pred = model_apertura.predict(X)
        y_pred_inverse = scaler_apertura.inverse_transform(y_pred)
        return jsonify({
            "last_date": data_file['fecha'].iloc[-1],
            "next_prediction": float(y_pred_inverse[0][0]),
            "last_values": datos.flatten()[-ventana:].tolist()
        })
    except Exception as e:
        return jsonify({"error": str(e)}), 400


@app.route('/predict_last_clausura', methods=['GET'])
def predict_last_clausura():
    """
    Endpoint para predecir el próximo valor basado en los últimos datos del CSV
    """
    try:
        # Cargar datos históricos
        data_file = pd.read_csv('./database/corr_bitcoin_diario_clausura.csv')
        datos = data_file['clausura'].to_numpy().reshape(-1, 1)
        # Tomar los últimos 'ventana' valores
        last_window = datos[-ventana:]
        last_window_scaled = scaler_clausura.transform(last_window)
        # Preparar para predicción
        X = last_window_scaled.T.reshape(1, ventana, 1)
        # Hacer predicción
        y_pred = model_clausura.predict(X)
        y_pred_inverse = scaler_clausura.inverse_transform(y_pred)
        return jsonify({
            "last_date": data_file['fecha'].iloc[-1],
            "next_prediction": float(y_pred_inverse[0][0]),
            "last_values": datos.flatten()[-ventana:].tolist()
        })
    except Exception as e:
        return jsonify({"error": str(e)}), 400


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
