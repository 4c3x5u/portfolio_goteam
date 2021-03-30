from rest_framework import serializers
from main.models import Team, User
from uuid import UUID


class RegisterSerializer(serializers.Serializer):
    username = serializers.CharField(
        min_length=5,
        max_length=35,
        error_messages={'required': 'Username cannot be empty.'}
    )
    password = serializers.CharField(
        min_length=8,
        max_length=255,
        error_messages={'required': 'Password cannot be empty.'}
    )
    password_confirmation = serializers.CharField(
        min_length=8,
        max_length=255,
        error_messages={'required': 'Password confirmation cannot be empty.'}
    )
    invite_code = serializers.CharField(required=False)
    is_admin = serializers.BooleanField()

    class Meta:
        model = User
        fields = ('username', 'password', 'password_confirmation', 'team',
                  'invite_code', 'is_admin')

    @staticmethod
    def validate_invite_code(value):
        if value:
            try:
                return UUID(value)
            except (ValueError, TypeError):
                raise serializers.ValidationError('Invalid invite code.')
        return value

    def create(self, validated_data):
        password = validated_data.get('password')
        password_confirmation = validated_data.get('password_confirmation')
        if password == password_confirmation:
            invite_code = validated_data.get('invite_code')
            if invite_code:
                try:
                    team = Team.objects.get(invite_code=invite_code)
                except Team.DoesNotExist:
                    raise serializers.ValidationError({
                        'invite_code': 'Team not found.'
                    })
                validated_data['team'] = team
                validated_data['is_admin'] = False
                validated_data.pop('invite_code')
            else:
                team = Team.objects.create()
                validated_data['team'] = team
                validated_data['is_admin'] = True
            validated_data.pop('password_confirmation')
            return User.objects.create(**validated_data)
        raise serializers.ValidationError({
            'password_confirmation': 'Confirmation does not match the '
                                     'password.'
        })

    def update(self, instance, validated_data):
        pass
