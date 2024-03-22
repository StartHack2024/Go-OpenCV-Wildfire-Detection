import cv2
import numpy as np
import requests
from datetime import datetime
from tensorflow.keras.models import load_model
import geocoder
import pytz

# Load the TensorFlow model
model = load_model('fine_tuned_model_EfficientNetB0.h5')

# Open the camera
cap = cv2.VideoCapture(0)

# Flag to prevent multiple POST requests
request_sent = False

while True:
    # Capture frame-by-frame
    ret, frame = cap.read()

    # Reshape the frame to add a batch dimension
    frame_batch = np.expand_dims(frame, axis=0)

    # Perform prediction on the frame
    prediction = model.predict(frame_batch)
    print(prediction[0])

    # If the prediction is 1 and no request has been sent, send POST request
    if np.round(prediction[0]) == 1 and not request_sent:
        # Obtain the current location
        g = geocoder.ip('me')
        
        # Prepare data for POST request
        data = {
            "datetime": str(datetime.now(pytz.utc)),
            "latitude": g.latlng[0] if g.latlng else 0.0,
            "longitude": g.latlng[1] if g.latlng else 0.0,
        }
        
        # Send POST request
        try:
            response = requests.post('http://localhost:8080/fireAlert', json=data)
            print("POST request sent successfully.")
            # Mark that the request has been sent
            request_sent = True
        except requests.exceptions.RequestException as e:
            print("Error sending POST request:", e)

    # Optional: Reset the request_sent flag under certain conditions here

    # Display the resulting frame
    cv2.imshow('frame', frame)

    # Break the loop if 'q' is pressed
    if cv2.waitKey(1) & 0xFF == ord('q'):
        break

# Release the capture
cap.release()
cv2.destroyAllWindows()