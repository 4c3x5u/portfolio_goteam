from rest_framework import serializers
from ..models import Task


class TaskSerializer(serializers.ModelSerializer):
    title = serializers.CharField( max_length=50, error_messages={
        'blank': 'Title cannot be empty.',
        'max_length': 'Title cannot be longer than 50 characters.'
    })

    class Meta:
        model = Task
        fields = '__all__'
