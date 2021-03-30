from rest_framework import serializers
from main.models import Team, User


class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = User
        fields = ('username', 'password', 'team')
