from rest_framework import serializers
from rest_framework.response import Response
from ..models import User


class LoginSerializer(serializers.Serializer):
    username = serializers.CharField(
        min_length=5,
        max_length=35,
        error_messages={'blank': 'Username cannot be empty.'}
    )
    password = serializers.CharField(min_length=8, max_length=255)

    def validate(self, data):
        if not data.get('username'):
            raise serializers.ValidationError({
                'username': 'Username cannot be empty.'
            }, 400)
        if not data.get('password'):
            raise serializers.ValidationError({
                'password': 'Password cannot be empty.'
            }, 400)
        user = User.objects.get(username=data.get('username'))
        if not user:
            raise serializers.ValidationError({
                'username': 'Invalid username.'
            }, 404)
        if user.password != data.get('password'):
            raise serializers.ValidationError({
                'password': 'Incorrect password.',
            }, 400)
        return super().validate(data)

    def create(self, validated_data):
        return Response({validated_data['username']: 'Login successful.'}, 200)

    def update(self, instance, validated_data):
        pass
