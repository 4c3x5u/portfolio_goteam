from rest_framework import serializers
from main.models import Team, User
from ..validation.val_custom import CustomAPIException
import bcrypt
import status


class UserSerializer(serializers.ModelSerializer):
    username = serializers.CharField(
        min_length=5,
        max_length=35,
        error_messages={
            'blank': 'Username cannot be empty.',
            'max_length': 'Username cannot be longer than 35 characters.'
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


class RegisterSerializer(UserSerializer):
    password_confirmation = serializers.CharField(
        error_messages={
            'blank': 'Password confirmation cannot be empty.'
        }
    )
    invite_code = serializers.UUIDField(
        required=False,
        error_messages={'invalid': 'Invalid invite code.'}
    )

    @staticmethod
    def validate_username(value):
        try:
            User.objects.get(username=value)
        except User.DoesNotExist:
            return value
        else:
            raise CustomAPIException('username',
                                     'Username already exists.',
                                     status.HTTP_400_BAD_REQUEST)

    def validate(self, attrs):
        password_confirmation = attrs.pop('password_confirmation')

        if password_confirmation != attrs.get('password'):
            raise CustomAPIException(
                'password_confirmation',
                'Confirmation must match the password.',
                status.HTTP_400_BAD_REQUEST
            )

        invite_code = attrs.get('invite_code')
        if invite_code:
            try:
                attrs['team'] = Team.objects.get(invite_code=invite_code)
            except Team.DoesNotExist:
                raise CustomAPIException('invite_code',
                                         'Team not found.',
                                         status.HTTP_400_BAD_REQUEST)
            attrs['is_admin'] = False
            attrs.pop('invite_code')
        else:
            attrs['is_admin'] = True

        return super().validate(attrs)

    def create(self, validated_data):
        if validated_data.get('is_admin') and not validated_data.get('team'):
            validated_data['team'] = Team.objects.create()

        validated_data['password'] = bcrypt.hashpw(
            bytes(validated_data['password'], 'utf-8'),
            bcrypt.gensalt()
        )

        return User.objects.create(**validated_data)

    def to_representation(self, instance):
        return {
            'msg': 'Registration successful.',
            'username': instance.username,
            'token': bcrypt.hashpw(
                bytes(instance.username, 'utf-8') + instance.password,
                bcrypt.gensalt()
            ).decode('utf-8'),
            'teamId': instance.team_id,
            'isAdmin': instance.is_admin
        }


class LoginSerializer(UserSerializer):
    def validate(self, attrs):
        try:
            user = User.objects.get(username=attrs.get('username'))
        except User.DoesNotExist:
            raise CustomAPIException('username',
                                     'Invalid username.',
                                     status.HTTP_400_BAD_REQUEST)

        pw_bytes = bytes(attrs.get('password'), 'utf-8')
        if not bcrypt.checkpw(pw_bytes, bytes(user.password)):
            raise CustomAPIException('password',
                                     'Invalid password.',
                                     status.HTTP_400_BAD_REQUEST)

        return user

    def to_representation(self, instance):
        return {
            'msg': 'Login successful.',
            'username': instance.username,
            'token': bcrypt.hashpw(
                bytes(instance.username, 'utf-8') + instance.password,
                bcrypt.gensalt()
            ).decode('utf-8'),
            'teamId': instance.team_id,
            'isAdmin': instance.is_admin,
        }