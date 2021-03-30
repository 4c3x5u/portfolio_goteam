from rest_framework import serializers
from rest_framework.response import Response
from main.models import Team, User
from uuid import UUID


class UserSerializer(serializers.Serializer):
    username = serializers.CharField(min_length=5, max_length=35)
    password = serializers.CharField(min_length=8, max_length=255)
    password_confirmation = serializers.CharField(min_length=8, max_length=255)
    invite_code = serializers.CharField(required=False, default='')

    class Meta:
        model = User
        fields = ('username', 'password', 'password_confirmation', 'team',
                  'is_admin', 'invite_code')

    def update(self, instance, validated_data):
        validated_data.pop('password_confirmation')
        validated_data.pop('invite_code')
        return User.objects.update(**validated_data)

    def validate_invite_code(self, value):
        try:
            return UUID(value)
        except (ValueError, TypeError):
            raise serializers.ValidationError('Invalid invite code.')

    def create(self, validated_data):
        password = validated_data.get('password')
        password_confirmation = validated_data.get('password_confirmation')
        if password == password_confirmation:
            try:
                invite_code = validated_data.get('invite_code')
            except KeyError:
                team = Team.objects.create()
                validated_data['team'] = team.id
                validated_data['is_admin'] = True
            else:
                team = Team.objects.get(invite_code=invite_code)
                validated_data['team'] = team.id
                validated_data['is_admin'] = False
                validated_data.pop('invite_code')
            validated_data.pop('password_confirmation')
            return User.objects.create(**validated_data)
        else:
            raise serializers.ValidationError('BABAAAN')
