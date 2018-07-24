import datetime

amount = float("1802.0")

a, b = divmod(amount, 100)
if b != 0:
    print(round((amount // 100 + 1) * 100, 1))
else:
    print(amount)


print(datetime.datetime.utcnow().isoformat())
