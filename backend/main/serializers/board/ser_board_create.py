from rest_framework import serializers

from main.serializers.board.ser_board import BoardSerializer
from main.validation.val_auth import authenticate, authorization_error, \
    authorize
from main.util import create_board


class CreateBoardSerializer(BoardSerializer):
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta(BoardSerializer.Meta):
        fields = 'auth_user', 'auth_token', 'team', 'name'

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
            raise board_serializer.errors

        return {'board_name': board_name,
                'team_id': team.id,
                'team_admin': team.user_set.get(username=user.username)}

    def create(self, validated_data):
        return create_board(name=validated_data.get('board_name'),
                            team_id=validated_data.get('team_id'),
                            team_admin=validated_data.get('team_admin'))

    def to_representation(self, instance):
        return {'msg': 'Board creation successful.',
                'id': instance.id}
