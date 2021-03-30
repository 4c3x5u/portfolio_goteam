from rest_framework import serializers
from main.models import Team, User


class UserSerializer(serializers.ModelSerializer):
    username = serializers.CharField(min_length=5, max_length=35)
    password = serializers.CharField(min_length=8, max_length=255)

    class Meta:
        model = User
        fields = ('username', 'password', 'team', 'is_admin')
