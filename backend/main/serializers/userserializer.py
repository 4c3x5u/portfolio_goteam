from rest_framework import serializers
from main.models import User


class UserSerializer(serializers.ModelSerializer):
    username = serializers.CharField(
        min_length=5,
        max_length=35,
        error_messages={
            'blank': 'Username cannot be empty.',
            'null': 'Username cannot be null.',
            'max_length': 'Username cannot be longer than 35 characters.',
            'does_not_exist': 'User not found'
        }
    )
    password = serializers.CharField(
        min_length=8,
        max_length=255,
        error_messages={
            'blank': 'Password cannot be empty.',
            'max_length': 'Password cannot be longer than 255 characters.'
        }
    )
    team = serializers.IntegerField(required=False)

    class Meta:
        model = User
        fields = '__all__'
