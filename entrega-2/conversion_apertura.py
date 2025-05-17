import tensorflow as tf

try:
    # 1. Cargar el modelo
    model = tf.keras.models.load_model('mejor_modelo_apertura.h5')

    # 2. Configurar el conversor
    converter = tf.lite.TFLiteConverter.from_keras_model(model)
    converter.optimizations = [tf.lite.Optimize.DEFAULT]
    converter.target_spec.supported_ops = [
        tf.lite.OpsSet.TFLITE_BUILTINS,
        tf.lite.OpsSet.SELECT_TF_OPS
    ]

    # 3. Convertir
    tflite_model = converter.convert()

    # 4. Guardar
    with open('modelo_apertura.tflite', 'wb') as f:
        f.write(tflite_model)
    print("¡Conversión exitosa!")

except Exception as e:
    print(f"Error durante la conversión: {str(e)}")
    if "Unsupported Ops" in str(e):
        print("--> Prueba agregar más opciones a target_spec.supported_ops")
