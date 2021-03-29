from rest_framework import serializers
from main.models import Team, User


class UserSerializer(serializers.ModelSerializer):
    password = serializers.CharField(max_length=255,
                                     style={'input_type': 'password'})
    class Meta:
        model = User
        fields = ('username', 'password', 'team')
