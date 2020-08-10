const osuReplayParser = require('osureplayparser');
const replayPath = "C:\\Users\\hndada\\Documents\\GitHub\\hndada\\gosu\\test\\MuangMuangE - Hideyuki Fukasawa - kengengreat three [Normal] (2020-07-30) OsuMania.osr";
const replay = osuReplayParser.parseReplay(replayPath);
console.log(replay["replay_data"][0]);