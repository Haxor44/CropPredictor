from rest_framework.response import Response
from rest_framework.decorators import api_view
import joblib
import json
import numpy as np

@api_view(['GET'])
def getData(request):
	person = {'name':'Toto','age':20}
	return Response(person)

@api_view(['POST'])
def getCrop(request):
	crop = joblib.load('cropPredictor.sav')
	data = json.loads(request.body)
	data2 = data[0]
	print([data2["measurements"]])
	predictedCrop = crop.predict([data2["measurements"]])	
	print(predictedCrop)
	return Response(predictedCrop)