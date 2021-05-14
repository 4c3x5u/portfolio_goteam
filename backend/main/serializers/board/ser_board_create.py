from rest_framework import serializers
from rest_framework.exceptions import ValidationError

from main.serializers.board.ser_board import BoardSerializer
from main.validation.val_auth import authenticate, authorize
from main.helpers import BoardHelper


class CreateBoardSerializer(BoardSerializer):
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta(BoardSerializer.Meta):
        fields = 'name', 'team', 'auth_user', 'auth_token'

    def validate(self, attrs):
        user = authenticate(attrs.get('auth_user'), attrs.get('auth_token'))

        team = attrs.get('team')
        authorize(user, team.id)

        board_name = attrs.get('name')
        board_serializer = BoardSerializer(data={
            'team': team.id,
            'name': board_name,
        })
        if not board_serializer.is_valid():
            raise ValidationError({'board': board_serializer.errors})

        return {'board_name': board_name,
                'team_id': team.id,
                'team_admin': team.user_set.get(username=user.username)}

    def create(self, validated_data):
        return BoardHelper.create(name=validated_data.get('board_name'),
                                  team_id=validated_data.get('team_id'),
                                  team_admin=validated_data.get('team_admin'))

    def to_representation(self, instance):
        return {'msg': 'Board creation successful.',
                'id': instance.id}
