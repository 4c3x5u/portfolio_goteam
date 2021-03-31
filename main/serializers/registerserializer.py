from rest_framework import serializers
from main.models import Team, User, Board

class RegisterSerializer(serializers.Serializer):
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
    password_confirmation = serializers.CharField(
        min_length=8,
        max_length=255,
        error_messages={'blank': 'Password confirmation cannot be empty.'}
    )
    invite_code = serializers.UUIDField(
        required=False,
        error_messages={'invalid': 'Invalid invite code.'}
    )
    is_admin = serializers.BooleanField()

    def validate(self, data):
        if data.get('password') != data.get('password_confirmation'):
            raise serializers.ValidationError({
                'password_confirmation': 'Confirmation does not match the '
                                         'password.'
            })
        invite_code = data.get('invite_code')
        if invite_code:
            try:
                team = Team.objects.get(invite_code=invite_code)
            except Team.DoesNotExist:
                raise serializers.ValidationError({
                    'invite_code': 'Team not found.'
                })
            data['team'] = team
            data['is_admin'] = False
        else:
            team = Team.objects.create()
            Board.objects.create(team=team)
            data['team'] = team
            data['is_admin'] = True
        return super().validate(data)

    def create(self, validated_data):
        if validated_data.get('invite_code'):
            validated_data.pop('invite_code')
        if validated_data.get('password_confirmation'):
            validated_data.pop('password_confirmation')
        return User.objects.create(**validated_data)

    def update(self, instance, validated_data):
        return super().update(instance, validated_data)
