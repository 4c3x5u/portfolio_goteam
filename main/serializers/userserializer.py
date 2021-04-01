from rest_framework import serializers
from main.models import Team, User, Board


class UserSerializer(serializers.ModelSerializer):
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
        required=False
    )
    team = serializers.PrimaryKeyRelatedField(queryset=Team.objects.all(),
                                              required=False)
    is_admin = serializers.BooleanField(default=True)
    invite_code = serializers.UUIDField(required=False, error_messages={
        'invalid': 'Invalid invite code.'
    })

    class Meta:
        model = User
        fields = ('username', 'password', 'password_confirmation', 'is_admin',
                  'team', 'invite_code')

    def validate(self, data):
        invite_code = data.get('invite_code')
        if invite_code:
            try:
                data['team'] = Team.objects.get(invite_code=invite_code)
            except Team.DoesNotExist:
                raise serializers.ValidationError({
                    'invite_code': 'Team not found.'
                })
            data['is_admin'] = False
            data.pop('invite_code')
        return super().validate(data)

    def create(self, validated_data):
        password_confirmation = validated_data.get('password_confirmation')
        if not password_confirmation:
            raise serializers.ValidationError({
                'password_confirmation': 'Password confirmation cannot be '
                                         'empty.'
            }, 'blank')
        if password_confirmation != validated_data.get('password'):
            raise serializers.ValidationError({
                'password_confirmation': 'Confirmation does not match the '
                                         'password.'
            }, 'no_match')
        validated_data.pop('password_confirmation')
        if validated_data.get('is_admin') and not validated_data.get('team'):
            team = Team.objects.create()
            Board.objects.create(team=team)
            validated_data['team'] = team
        return User.objects.create(**validated_data)
