CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS photo_locations (
    location_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    geom geometry(Point, 4326) NOT NULL,
    geohash VARCHAR(9) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_photo_locations_geom_gist ON photo_locations USING GIST (geom);
CREATE INDEX IF NOT EXISTS idx_photo_locations_geohash ON photo_locations (geohash);
CREATE INDEX IF NOT EXISTS idx_photo_locations_user_id ON photo_locations(user_id);

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS tips(
    tips_id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO tips (tips_id, title, content) VALUES
(gen_random_uuid(),'ここはなに？','やくに立つかもしれない情報が書かれているよ！'),
(gen_random_uuid(),'TIPS','これはチップスって読むんだよ！'),
(gen_random_uuid(),'TIPS','これはティップスって読むんだよ！'),
(gen_random_uuid(),'TIPS','これはテップスって読むんだよ！'),
(gen_random_uuid(),'TIPS','これはタイプスって読むんだよ！'),
(gen_random_uuid(),'制作者１','「開発疲れました(泣)」'),
(gen_random_uuid(),'制作者？２','「九割九分九厘何もしていない」'),
(gen_random_uuid(),'制作者３','「制作者２は何をしていたんだ...」'),
(gen_random_uuid(),'制作者４','「彼女募集中」'),
(gen_random_uuid(),'制作者５','「HipsはHipの複数形」'),
(gen_random_uuid(),'制作者６','「Tipsってチップスじゃないらしい」'),
(gen_random_uuid(),'ヒートマップ','ヒートマップの位置情報はだいたいだから、きれいな景色を歩いて探してみてね！'),
(gen_random_uuid(),'サーバー','サーバーがとても貧弱なので強化投げ銭ください！！'),
(gen_random_uuid(),'三分割法','画面を縦横3分割し、交点に被写体を配置するとバランスの良い写真になるよ。'),
(gen_random_uuid(),'日の丸構図','被写体を中央に配置するシンプルな構図。被写体の存在感を強く出したいときに有効だよ。'),
(gen_random_uuid(),'対角線構図','画面を斜めに使う構図。動きや奥行きを表現しやすいよ。'),
(gen_random_uuid(),'放射線構図','線が一点に集まるように配置すると、自然と視線を誘導できるよ。'),
(gen_random_uuid(),'額縁構図','窓や木などで被写体を囲むと、写真に奥行きと物語性が出るよ。'),
(gen_random_uuid(),'シンメトリー','左右や上下を対称にすると、整った印象や美しさを強調できるよ。'),
(gen_random_uuid(),'前ボケ','手前に物を入れてぼかすことで、立体感や雰囲気を演出できるよ。'),
(gen_random_uuid(),'後ボケ','背景をぼかして被写体を引き立てる基本テクニックだよ。'),
(gen_random_uuid(),'余白を活かす','被写体の周囲に余白を残すと、写真に落ち着きや印象的な雰囲気が生まれるよ。'),
(gen_random_uuid(),'ローアングル','下から撮ることで、被写体を大きく力強く見せられるよ。'),
(gen_random_uuid(),'ハイアングル','上から撮ると、被写体を小さく可愛らしく見せられるよ。'),
(gen_random_uuid(),'S字構図','道や川などをS字に入れると、自然な視線誘導ができるよ。'),
(gen_random_uuid(),'C字構図','被写体を包み込むような曲線で、柔らかい印象になるよ。'),
(gen_random_uuid(),'フレームアウト','あえて被写体の一部を切ることで、臨場感を出せるよ。'),
(gen_random_uuid(),'反復構図','同じ形や色を繰り返すと、リズム感のある写真になるよ。'),
(gen_random_uuid(),'一点集中','他の情報を減らし、主役だけを際立たせる構図だよ。'),
(gen_random_uuid(),'奥行きを意識','前景・中景・背景を意識すると、立体的な写真になるよ。'),
(gen_random_uuid(),'視線の先に余白','人物や動物の視線の先を空けると、自然で印象的になるよ。'),
(gen_random_uuid(),'水平・垂直を意識','建物や風景では線をまっすぐにすると安定感が出るよ。'),
(gen_random_uuid(),'色の対比','補色関係の色を使うと、被写体が際立つよ。'),
(gen_random_uuid(),'光を構図に入れる','逆光や木漏れ日を使うと、ドラマチックな写真になるよ。'),
(gen_random_uuid(),'影を主役にする','影を意識すると、抽象的で印象的な写真が撮れるよ。'),
(gen_random_uuid(),'空を大きく入れる','空の割合を増やすと、開放感のある写真になるよ。'),
(gen_random_uuid(),'あえて傾ける','水平を崩すことで、緊張感や動きを表現できるよ。'),
(gen_random_uuid(),'被写体を端に寄せる','中央を外すことで、写真に余韻やストーリーが生まれるよ'),
(gen_random_uuid(),'最後の一押し！','最後は自分を信じるのだ！！'),
(gen_random_uuid(),'共有','撮影後は各SNSに簡単に共有できます！(宣伝してね(嘘))'),
(gen_random_uuid(),'雑談','開発者の一人「夜景を見に行きたい」'),
(gen_random_uuid(),'雑談2','バグを見つけたら暖かい目で見てください...'),
(gen_random_uuid(),'開発者の裏側','(ノーマルカメラでは)物撮るってレベルじゃねぇぞ！'),
(gen_random_uuid(),'開発者の裏側','「限界超えてんだよ！」'),
(gen_random_uuid(),'ヒートマップ3','ヒートマップは色が濃いほうがたくさん取られています！'),
(gen_random_uuid(),'カニ歩き撮影','カメラを右や左に動かすときは、スマホの向きを変えるのではなく、自分自身が横に並行移動するのがコツだよ！'),
(gen_random_uuid(),'ズームより一歩','画角を調整したいときは画面をピンチするより、自分自身が一歩前へ進むか後ろへ下がるほうが画質が綺麗に保てるよ！'),
(gen_random_uuid(),'スマホ逆さま撮影','地面に近い被写体を撮るときは、スマホを上下逆さまに持ってレンズを地面に近づけると、驚くほど迫力が出るよ。'),
(gen_random_uuid(),'料理は「回す」','お皿の向きを少しずつ変えてみて、影が一番綺麗に落ちる角度、具材が一番おいしそうに見える「顔」を探そう。'),
(gen_random_uuid(),'背景から離れる','被写体を背景から少し離して置く（あるいは離れて立ってもらう）と、後ろが自然にボケて主役が浮き上がるよ。'),
(gen_random_uuid(),'肘を固定！','スマホを動かす指示が出たとき、脇を締めて肘を体に固定しながら動くと、手ブレを最小限に抑えられるよ。'),
(gen_random_uuid(),'光源の位置を確認','光が背中側にある「順光」だと色が綺麗に、横から当たる「斜光」だと凹凸がはっきりして立体感が出るよ。'),
(gen_random_uuid(),'引き算の美学','画面の中に色んなものを入れすぎないこと。主役以外を画面の外に出す勇気が、良い写真への近道だよ。'),
(gen_random_uuid(),'視線の抜け感','人物や動物を撮るときは、その視線が向いている方向に広いスペース（余白）を作ると、開放的な写真になるよ。'),
(gen_random_uuid(),'リズムとパターン','同じ形のものが並んでいる場所では、その「繰り返し」を画面いっぱいに収めると、デザイン性の高い一枚になるよ。'),
(gen_random_uuid(),'物語を想像させる','あえて全体を映さず一部を切り取ると、見る人が「その先」を想像して、写真にストーリー性が生まれるよ。'),
(gen_random_uuid(),'黄金比の魔法','三分割法よりも少しだけ中央に寄せた「黄金比」を意識すると、より自然で洗練された安定感が手に入るよ。'),
(gen_random_uuid(),'開発者のボヤキ','「カメラを右に」って指示が出たのに左に行っちゃう人、僕も仲間です。'),
(gen_random_uuid(),'開発者の裏側','実はAIの構図アドバイス、たまに僕らよりセンスが良いときがあって凹みます。'),
(gen_random_uuid(),'筋肉痛','低い位置から撮りすぎて翌日筋肉痛になりました。みんなは無理しないでね。'),
(gen_random_uuid(),'Tipsの極意','結局のところ、たくさんシャッターを切った人が一番上達するんだって。開発者もコード書きます。');
