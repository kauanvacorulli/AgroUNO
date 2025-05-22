const int sensorPin = A0; // Defines sensor pin to analog input A0

void setup() {
  Serial.begin(9600);
}

void loop() {
  int sensorValue = analogRead(sensorPin);
  
  // Convert the raw value from 0 to 1023 to 0 to 100, with 0% being dry and 100% being soaking wet
  int moisturePercent = map(sensorValue, 1023, 0, 0, 100);

  Serial.println(moisturePercent); // Prints percentage
  delay(1000); // Delays loop by a second
}
