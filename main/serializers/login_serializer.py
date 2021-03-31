from rest_framework import serializers
from main.models import User


class LoginSerializer(serializers.Serializer):
    username = serializers.CharField(
        min_length=5,
        max_length=35,
        error_messages={'blank': 'Username cannot be empty.'}
    )
    password = serializers.CharField(
        min_length=8,
        max_length=255,
        error_messages={'blank': 'Password cannot be empty.'}
    )

    def validate(self, data):
        user = User.objects.get(username=data.get('username'))
        if not user:
            raise serializers.ValidationError({
                'username': 'Invalid username.'
            })
        if user.password != data.get('password'):
            raise serializers.ValidationError({
                'password': 'Invalid password.',
            })
        return super().validate(data)

    def create(self, validated_data):
        pass

    def update(self, instance, validated_data):
        pass
