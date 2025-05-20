import sys
import joblib
import seaborn as sns
import pandas as pd
import numpy as np
from sklearn.metrics import mean_squared_error
from sklearn.preprocessing import MinMaxScaler
from keras.layers import Dense, GRU
from keras import Sequential
import tensorflow as tf
import logging
logging.getLogger("tensorflow").setLevel(logging.ERROR)
tf.autograph.set_verbosity(0)


# funciones auxiliares
def calcular_rmse(y, y_hat):
    return np.sqrt(mean_squared_error(y, y_hat))


def generar_ventanas(datos_preproc, ventana):
    X, y = [], []
    for i in range(len(datos_preproc) - ventana):
        X.append(datos_preproc[i:i+ventana, 0])
        y.append(datos_preproc[i+ventana, 0])
    X, y = np.array(X), np.array(y)
    return X, y


data_file = pd.read_csv("../database/bitcoin_apertura_2013_actual.csv")
data_file.head()
data_file_apertura = data_file["apertura"]
data_file_apertura.head()

# Conversión a array de numpy
datos = data_file_apertura.to_numpy().reshape(-1, 1)
print("Dimensiones del conjunto de datos: ", datos.shape)

# Estandarización min-max
scaler = MinMaxScaler()
datos_estandarizados = scaler.fit_transform(datos)
datos_estandarizados[:10]

# generacion de conjuntos de prueba y entrenamiento
ventana = 8  # Tomaremos los 7 valores previos
n = len(datos)
m = n-ventana
tam_entrenamiento = int(n*0.7)
tam_prueba = n-tam_entrenamiento-ventana
print("Vectores para entrenamiento:", tam_entrenamiento)
print("Vectores para prueba: ", tam_prueba)

# Generamos vectores de características.
# Cada vector consiste en el precio x y los V precios anteriores,
# donde V es el tamaño de la ventana
# En este caso, la etiqueta numérica sería el precio x
X, y = generar_ventanas(datos_estandarizados, ventana)
# Preparamos tensores para poder utilizar GRU
X = X.reshape(X.shape[0], X.shape[1], 1)
X[:3]
y[:3]

# Genera X_train, X_test, y_train, y y_test
# a partir de X, y y tam_entrenamiento.
X_train, X_test = X[:tam_entrenamiento], X[tam_entrenamiento:]
y_train, y_test = y[:tam_entrenamiento], y[tam_entrenamiento:]
y_train = y_train.reshape(-1, 1)
y_test = y_test.reshape(-1, 1)
print("Dimensiones de X: ", X.shape)
print("Dimensiones de y: ", y.shape)
print("Dimensiones de X_train: ", X_train.shape)
print("Dimensiones de X_test: ", X_test.shape)
print("Dimensiones de y_train: ", y_train.shape)
print("Dimensiones de y_test: ", y_test.shape)

# Generación de modelos de aprendizaje
modelos = []
modelo_a = Sequential()
modelo_a.add(GRU(units=50, return_sequences=False, input_shape=(ventana, 1)))
modelo_a.add(Dense(1))
modelos.append(modelo_a)
nombres_modelos = ["Modelo A"]


@tf.autograph.experimental.do_not_convert
def run(modelos, nombres_modelos, X_train, X_test, y_train, y_test, scaler):
    # Listas para guardar las predicciones de entrenamiento, las predicciones
    # de prueba y los errores
    Y_train_pred_inverse = []
    Y_test_pred_inverse = []
    errores = []
    mejor_mse = sys.float_info.max
    corridas = 5
    for k in range(corridas):
        print("--------------Probando ", nombres_modelos[0], "---------------")
        modelo = modelos[0]
        modelo.compile(loss="mean_squared_error", optimizer="adam")
        modelo.fit(X_train, y_train, epochs=5, batch_size=2, verbose=1)
        y_train_pred = modelo.predict(X_train, verbose=1)
        y_test_pred = modelo.predict(X_test, verbose=1)
        rmse_ejecucion = calcular_rmse(y_test, y_test_pred)
        print("   Raíz del error medio cuadrático (RMSE):",
              round(rmse_ejecucion, 4))
        errores.append(rmse_ejecucion)
        if rmse_ejecucion < mejor_mse:
            mejor_mse = rmse_ejecucion
            mejor_modelo = modelo
    print("Mejor MSE: ", mejor_mse)
    # errores.append(errs_modelo)
    # Se vuelven a calcular con el mejor modelo (también se pueden guardar)
    y_train_pred = mejor_modelo.predict(X_train, verbose=1)
    y_test_pred = mejor_modelo.predict(X_test, verbose=1)
    y_train_pred_inverse = scaler.inverse_transform(y_train_pred)
    y_test_pred_inverse = scaler.inverse_transform(y_test_pred)
    Y_train_pred_inverse.append(y_train_pred_inverse)
    Y_test_pred_inverse.append(y_test_pred_inverse)

    # !!! Ejercicio 5
    # Guarda el modelo en formato HDF5
    mejor_modelo.save("mejor_modelo_apertura.h5")
    # Guarda el scaler para poder preprocesar nuevos datos
    joblib.dump(scaler, "scaler_apertura.save")
    return [Y_train_pred_inverse, Y_test_pred_inverse, errores]


[Y_train_pred_inverse, Y_test_pred_inverse, errores] = run(
    modelos, nombres_modelos, X_train, X_test, y_train, y_test, scaler)

print(modelos[0].summary())
print()

# errores de modelo
datos_errores = {}
datos_errores[nombres_modelos[0]] = errores
df_errores = pd.DataFrame(datos_errores)
print(errores)
