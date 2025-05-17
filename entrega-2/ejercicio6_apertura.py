from sklearn.metrics import mean_squared_error
import pandas as pd
import numpy as np
from keras.models import load_model
from sklearn.preprocessing import MinMaxScaler
import joblib
import matplotlib.pyplot as plt

# 1. Cargar los datos originales
data_file = pd.read_csv('corr_bitcoin_diario_apertura.csv')
data_file_apertura = data_file['apertura']
datos = data_file_apertura.to_numpy().reshape(-1, 1)

# 2. Cargar el scaler guardado
scaler = joblib.load('scaler_apertura.save')

# 3. Cargar el modelo
# Usa el nombre que guardaste
modelo_cargado = load_model('mejor_modelo_apertura.h5')

# 4. Preprocesamiento (debe ser IDÉNTICO al que
# usaste durante el entrenamiento)
datos_estandarizados = scaler.transform(
    datos)  # Usa transform() no fit_transform()

# 5. Recrear las ventanas temporales (mismo tamaño
# que durante el entrenamiento)
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

# 8. Visualización
plt.figure(figsize=(12, 6))
plt.plot(y_real_inverse, label='Precio Real', color='blue', alpha=0.6)
plt.plot(y_pred_inverse, label='Predicciones', color='red', alpha=0.6)
plt.title('Comparación: Precio Real vs Predicciones del Modelo Cargado')
plt.xlabel('Tiempo')
plt.ylabel('Precio (USD)')
plt.legend()
plt.grid(True)
plt.show()

# 9. Calcular métricas de error (opcional)
rmse = np.sqrt(mean_squared_error(y_real_inverse, y_pred_inverse))
print(f'RMSE usando el modelo cargado: {rmse:.2f}')
