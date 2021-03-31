from rest_framework import serializers
from ..models import User, Team, Board


class CreateBoardSerializer(serializers.Serializer):
    username: serializers.CharField(min_length=5, max_length=35)
    team_id: serializers.IntegerField()

    def validate(self, data):
        try:
            user = User.objects.get(username=data.username)
        except User.DoesNotExist:
            raise serializers.ValidationError({
                'username': 'Invalid username.'
            })
        if not user.is_admin:
            raise serializers.ValidationError({
                'user_not_admin': 'Only the team admin can create a board.'
            })
        try:
            Team.objects.get(id=data.team_id)
        except Team.DoesNotExist:
            raise serializers.ValidationError({
                'team_id': 'Invalid team ID.'
            })
        return super().validate(data)

    def create(self, validated_data):
        Board.objects.create(team=validated_data.team_id)

    def update(self, instance, validated_data):
        pass

