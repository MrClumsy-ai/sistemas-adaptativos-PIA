import pandas as pd
import numpy as np
from keras.models import load_model
import joblib

# 1. Cargar los datos originales
data_file = pd.read_csv('corr_bitcoin_diario_clausura.csv')
data_file_clausura = data_file['clausura']
datos = data_file_clausura.to_numpy().reshape(-1, 1)

# 2. Cargar el scaler guardado
scaler = joblib.load('scaler_clausura.save')

# 3. Cargar el modelo
modelo_cargado = load_model('mejor_modelo_clausura.h5')

# 4. Preprocesamiento
datos_estandarizados = scaler.transform(datos)

# 5. Recrear las ventanas temporales
ventana = 8  # Debe ser el mismo valor usado originalmente


def generar_ventanas(datos_preproc, ventana):
    X, y = [], []
    for i in range(len(datos_preproc) - ventana):
        X.append(datos_preproc[i:i+ventana, 0])
        y.append(datos_preproc[i+ventana, 0])
    return np.array(X), np.array(y)


X, y = generar_ventanas(datos_estandarizados, ventana)
X = X.reshape(X.shape[0], X.shape[1], 1)  # Reformar para GRU

# 6. Hacer predicciones
y_pred = modelo_cargado.predict(X)

# 7. Invertir la normalización
y_pred_inverse = scaler.inverse_transform(y_pred)
y_real_inverse = scaler.inverse_transform(y.reshape(-1, 1))
