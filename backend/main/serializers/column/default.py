from rest_framework.serializers import ModelSerializer
from main.models import Column


class ColumnSerializer(ModelSerializer):
    class Meta:
        model = Column
        fields = '__all__'
