from rest_framework import serializers
from django.core.serializers import serialize
from ..models import Team, Board


class BoardSerializer(serializers.ModelSerializer):
    class Meta:
        model = Board
        fields = ('id',)


class ListBoardsSerializer(serializers.Serializer):
    team_id = serializers.IntegerField()

    def update(self, instance, validated_data):
        raise Exception('Update is not allowed in ListBoardsSerializer.')

    def create(self, validated_data):
        raise Exception('Create is not allowed in ListBoardsSerializer.')

    def validate(self, data):
        team_id = data.get('team_id')
        if not team_id:
            raise serializers.ValidationError({
                'team_id': 'Team ID cannot be empty.'
            })
        try:
            team = Team.objects.get(id=team_id)
        except Team.DoesNotExist:
            raise serializers.ValidationError({
                'team_id': 'Invalid team ID.',
            })
        if not Board.objects.filter(team=team):
            raise serializers.ValidationError(
                'No boards found for this team.',
                'not_found'
            )
        return super().validate(data)

    @staticmethod
    def get_list(team_id):
        boards = Board.objects.filter(team_id=team_id)
        boards_list = list(
            map(lambda board_id: {'board_id': board_id}, boards)
        )
        serializer = BoardSerializer(boards_list, many=True)
        return serializer.data
