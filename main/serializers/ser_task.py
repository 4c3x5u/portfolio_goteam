from rest_framework import serializers
from ..models import Task


class TaskSerializer(serializers.ModelSerializer):
    title = serializers.CharField(
        max_length=50,
        error_messages={
            'blank': 'Title cannot be empty.'
        }
    )

    class Meta:
        model = Task
        fields = '__all__'
